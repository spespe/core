// Rendezvous protocol implementation. The other name - Bidirectional Locator.
//
// Rendezvous is a protocol aimed to achieve mutual address resolution for both
// client and server. This is especially useful where there is no guarantee
// that both peers are reachable directly, i.e. behind a NAT for example.
// The protocol allows for servers to publish their private network addresses
// while resolving their real remove address. This information is saved under
// some ID until a connection held and heartbeat frames transmitted.
// When a client wants to perform a connection it goes to this server, informs
// about its own private network addresses and starting from this point a
// rendezvous can be achieved.
// Both client and server are informed about remote peers public and private
// addresses and they may perform several attempts to create an actual p2p
// connection.
//
// By trying to connect to private addresses they can reach each other if and
// only if they're both in the same LAN.
// If a peer doesn't have private address, the connection is almost guaranteed
// to be established directly from the private peer to the public one.
// If both of them are located under different NATs, a TCP punching attempt can
// be performed.
// At last, when there is no hope, a special relay server can be used to
// forward the traffic.
//
// Servers should publish all of possibly reachable endpoints for all protocols
// they support.
// Clients should specify the desired protocol and ID for resolution.
//
// Currently only TCP endpoints exchanging is supported.
// TODO: When resolving it's necessary to track also IP version. For example to be able not to return IPv6 when connecting socket is IPv4.

package rendezvous

import (
	"errors"
	"math/rand"
	"net"
	"sync"

	"github.com/sonm-io/core/insonmnia/auth"
	"github.com/sonm-io/core/proto"
	"github.com/sonm-io/core/util/xgrpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type deleter func()

type peerCandidate struct {
	Peer
	C chan<- Peer
}

type meeting struct {
	mu sync.Mutex
	// We allow multiple clients to be waited for servers.
	clients map[PeerID]peerCandidate
	// Also we allow the opposite: multiple servers can be registered for
	// fault tolerance.
	servers map[PeerID]peerCandidate
}

func newMeeting() *meeting {
	return &meeting{
		servers: map[PeerID]peerCandidate{},
		clients: map[PeerID]peerCandidate{},
	}
}

func (m *meeting) addServer(peer Peer, c chan<- Peer) {
	m.servers[peer.ID] = peerCandidate{Peer: peer, C: c}
}

func (m *meeting) addClient(peer Peer, c chan<- Peer) {
	m.clients[peer.ID] = peerCandidate{Peer: peer, C: c}
}

func (m *meeting) randomServer() *peerCandidate {
	if len(m.servers) == 0 {
		return nil
	}

	var keys []PeerID
	for key := range m.servers {
		keys = append(keys, key)
	}

	v := m.servers[keys[rand.Intn(len(keys))]]
	return &v
}

func (m *meeting) randomClient() *peerCandidate {
	if len(m.clients) == 0 {
		return nil
	}
	var keys []PeerID
	for key := range m.clients {
		keys = append(keys, key)
	}

	v := m.clients[keys[rand.Intn(len(keys))]]
	return &v
}

// Server represents a rendezvous server.
//
// This server is responsible for tracking servers and clients to make them
// meet each other.
type Server struct {
	cfg      Config
	log      *zap.Logger
	server   *grpc.Server
	resolver resolver

	mu sync.Mutex
	rv map[string]*meeting
}

// NewServer constructs a new rendezvous server using specified config and
// options.
//
// The server supports TLS by passing transport credentials using
// WithCredentials option.
// Also it is possible to activate logging system by passing a logger using
// WithLogger function as an option.
func NewServer(cfg Config, options ...Option) (*Server, error) {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	server := &Server{
		cfg: cfg,
		log: opts.log,
		server: xgrpc.NewServer(
			opts.log,
			xgrpc.Credentials(opts.credentials),
			xgrpc.DefaultTraceInterceptor(),
		),
		resolver: newExternalResolver(""),
		rv:       map[string]*meeting{},
	}

	server.log.Debug("configured authentication settings", zap.Any("credentials", opts.credentials.Info()))

	sonm.RegisterRendezvousServer(server.server, server)
	server.log.Debug("registered gRPC server")

	return server, nil
}

func (m *Server) Resolve(ctx context.Context, request *sonm.ConnectRequest) (*sonm.RendezvousReply, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errNoPeerInfo()
	}

	m.log.Info("resolving remote peer", zap.String("id", request.ID))

	id := request.ID
	peerHandle := NewPeer(*peerInfo, request.PrivateAddrs)

	c, deleter := m.addServerWatch(id, peerHandle)
	defer deleter()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case p := <-c:
		m.log.Info("providing remote server endpoint(s)",
			zap.String("id", request.ID),
			zap.Stringer("public_addr", p.Addr),
			zap.Any("private_addrs", p.privateAddrs),
		)
		return m.newReply(p)
	}
}

func (m *Server) Publish(ctx context.Context, request *sonm.PublishRequest) (*sonm.RendezvousReply, error) {
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errNoPeerInfo()
	}

	ethAddr, err := auth.ExtractWalletFromContext(ctx)
	if err != nil {
		return nil, err
	}

	m.log.Info("publishing remote peer", zap.String("id", ethAddr.String()))

	id := ethAddr.String()
	peerHandle := NewPeer(*peerInfo, request.PrivateAddrs)

	c, deleter := m.newClientWatch(id, peerHandle)
	defer deleter()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case p := <-c:
		m.log.Info("providing remote client endpoint(s)",
			zap.String("id", ethAddr.String()),
			zap.Stringer("public_addr", p.Addr),
			zap.Any("private_addrs", p.privateAddrs),
		)
		return m.newReply(p)
	}
}

func (m *Server) addServerWatch(id string, peer Peer) (<-chan Peer, deleter) {
	c := make(chan Peer, 1)

	m.mu.Lock()
	defer m.mu.Unlock()

	meeting, ok := m.rv[id]
	if ok {
		// Notify both sides immediately if there is match between candidates.
		if server := meeting.randomServer(); server != nil {
			c <- server.Peer
			server.C <- peer
		} else {
			meeting.addClient(peer, c)
		}
	} else {
		meeting := newMeeting()
		meeting.addClient(peer, c)
		m.rv[id] = meeting
	}

	return c, func() { m.removeServerWatch(id, peer) }
}

func (m *Server) newClientWatch(id string, peer Peer) (<-chan Peer, deleter) {
	c := make(chan Peer, 1)

	m.mu.Lock()
	defer m.mu.Unlock()

	meeting, ok := m.rv[id]
	if ok {
		if client := meeting.randomClient(); client != nil {
			c <- client.Peer
			client.C <- peer
		} else {
			meeting.addServer(peer, c)
		}
	} else {
		meeting := newMeeting()
		meeting.addServer(peer, c)
		m.rv[id] = meeting
	}

	return c, func() { m.removeClientWatch(id, peer) }
}

func (m *Server) removeClientWatch(id string, peer Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	candidates, ok := m.rv[id]
	if ok {
		delete(candidates.servers, peer.ID)
	}
}

func (m *Server) removeServerWatch(id string, peer Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	candidates, ok := m.rv[id]
	if ok {
		delete(candidates.clients, peer.ID)
	}
}

func (m *Server) newReply(peer Peer) (*sonm.RendezvousReply, error) {
	addr, err := sonm.NewAddr(peer.Addr)
	if err != nil {
		return nil, err
	}

	// If peer's public IP address suddenly arrived as a private this only
	// means that both Rendezvous and server/client are located within the
	// same VLAN. This should rarely happen, but in these cases we must
	// resolve public IP of ours.
	if addr.IsPrivate() {
		publicIP, err := m.resolver.PublicIP()
		if err != nil {
			return nil, err
		}

		addr.Addr.Addr = publicIP.String()
	}

	return &sonm.RendezvousReply{
		PublicAddr:   addr,
		PrivateAddrs: peer.privateAddrs,
	}, nil
}

func (m *Server) Info(ctx context.Context, request *sonm.Empty) (*sonm.RendezvousState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	state := make(map[string]*sonm.RendezvousMeeting)

	for id, meeting := range m.rv {
		clients := make(map[string]*sonm.RendezvousReply)
		servers := make(map[string]*sonm.RendezvousReply)

		for clientID, candidate := range meeting.clients {
			reply, err := m.newReply(candidate.Peer)
			if err != nil {
				return nil, err
			}

			clients[clientID.String()] = reply
		}

		for serverID, candidate := range meeting.servers {
			reply, err := m.newReply(candidate.Peer)
			if err != nil {
				return nil, err
			}
			servers[serverID.String()] = reply
		}

		state[id] = &sonm.RendezvousMeeting{
			Servers: servers,
			Clients: clients,
		}
	}

	return &sonm.RendezvousState{
		State: state,
	}, nil
}

// Run starts accepting incoming connections, serving them by blocking the
// caller execution context until either explicitly terminated using Stop
// or some critical error occurred.
//
// Always returns non-nil error.
func (m *Server) Run() error {
	listener, err := net.Listen(m.cfg.Addr.Network(), m.cfg.Addr.String())
	if err != nil {
		return err
	}

	m.log.Info("rendezvous is ready to serve", zap.Stringer("endpoint", listener.Addr()))
	return m.server.Serve(listener)
}

// Stop stops the server.
//
// It immediately closes all open connections and listeners. Also it cancels
// all active RPCs on the server side and the corresponding pending RPCs on
// the client side will get notified by connection errors.
func (m *Server) Stop() {
	m.log.Info("rendezvous is shutting down")
	m.server.Stop()
}

func errNoPeerInfo() error {
	return errors.New("no peer info provided")
}
