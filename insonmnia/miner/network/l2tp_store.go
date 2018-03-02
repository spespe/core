package network

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

type L2TPNetworkStore struct {
	mu       sync.Mutex
	aliases  map[string]string
	networks map[string]*L2TPNetwork
}

func NewL2TPNetworkStore() *L2TPNetworkStore {
	return &L2TPNetworkStore{
		aliases:  make(map[string]string),
		networks: make(map[string]*L2TPNetwork),
	}
}

func (s *L2TPNetworkStore) AddNetwork(netID string, netInfo *L2TPNetwork) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.aliases[netID]; ok {
		return errors.Errorf("network already exists: %s", netID)
	}

	s.aliases[netID] = netID
	s.networks[netID] = netInfo

	return nil
}

func (s *L2TPNetworkStore) AddNetworkAlias(netID, alias string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.aliases[netID]; !ok {
		return errors.Errorf("network not found: %s", netID)
	}

	s.aliases[alias] = netID

	return nil
}

func (s *L2TPNetworkStore) RemoveNetwork(netID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	translatedID, ok := s.aliases[netID]
	if !ok {
		return errors.Errorf("network not found: %s", netID)
	}

	delete(s.networks, translatedID)

	for alias, target := range s.aliases {
		if target == translatedID {
			delete(s.aliases, alias)
		}
	}

	return nil
}

func (s *L2TPNetworkStore) GetNetwork(netID string) (*L2TPNetwork, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	translatedID, ok := s.aliases[netID]
	if !ok {
		return nil, errors.Errorf("network not found: %s", netID)
	}

	if netInfo, ok := s.networks[translatedID]; ok {
		return netInfo, nil
	}

	return nil, errors.Errorf("network not found: %s", netID)
}

func (s *L2TPNetworkStore) GetNetworks() []*L2TPNetwork {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []*L2TPNetwork
	for _, netInfo := range s.networks {
		out = append(out, netInfo)
	}

	return out
}

type L2TPNetwork struct {
	ID          string
	PoolID      string
	count       int
	networkOpts *L2TPNetworkConfig
	store       *L2TPEndpointStore
}

func newNetworkInfo(opts *L2TPNetworkConfig) *L2TPNetwork {
	return &L2TPNetwork{
		networkOpts: opts,
		store:       NewL2TPEndpointStore(),
	}
}

func (n *L2TPNetwork) Setup() error {
	n.PoolID = n.networkOpts.GetHash()

	return nil
}

func (n *L2TPNetwork) ConnInc() {
	n.count++
}

type L2TPEndpointStore struct {
	mu        sync.Mutex
	aliases   map[string]string
	endpoints map[string]*L2TPEndpoint
}

func NewL2TPEndpointStore() *L2TPEndpointStore {
	return &L2TPEndpointStore{
		aliases:   make(map[string]string),
		endpoints: make(map[string]*L2TPEndpoint),
	}
}

func (s *L2TPEndpointStore) AddEndpoint(endpointID string, eptInfo *L2TPEndpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.aliases[endpointID]; ok {
		return errors.Errorf("endpoint already exists: %s", endpointID)
	}

	s.aliases[endpointID] = endpointID
	s.endpoints[endpointID] = eptInfo

	return nil
}

func (s *L2TPEndpointStore) AddEndpointAlias(endpointID, alias string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.aliases[endpointID]; !ok {
		return errors.Errorf("network not found: %s", endpointID)
	}

	s.aliases[alias] = endpointID

	return nil
}

func (s *L2TPEndpointStore) RemoveEndpoint(endpointID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	translatedID, ok := s.aliases[endpointID]
	if !ok {
		return errors.Errorf("endpoint not found: %s", endpointID)
	}

	delete(s.endpoints, translatedID)

	for alias, target := range s.aliases {
		if target == translatedID {
			delete(s.aliases, alias)
		}
	}

	return nil
}

func (s *L2TPEndpointStore) GetEndpoint(netID string) (*L2TPEndpoint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	translatedID, ok := s.aliases[netID]
	if !ok {
		return nil, errors.Errorf("endpoint not found: %s", netID)
	}

	if netInfo, ok := s.endpoints[translatedID]; ok {
		return netInfo, nil
	}

	return nil, errors.Errorf("endpoint not found: %s", netID)
}

func (s *L2TPEndpointStore) GetEndpoints() []*L2TPEndpoint {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []*L2TPEndpoint
	for _, eptInfo := range s.endpoints {
		out = append(out, eptInfo)
	}

	return out
}

type L2TPEndpoint struct {
	ID           string
	Name         string
	ConnName     string
	PPPOptFile   string
	PPPDevName   string
	AssignedCIDR string
	AssignedIP   string
	networkOpts  *L2TPNetworkConfig
}

func NewL2TPEndpoint(netInfo *L2TPNetwork) *L2TPEndpoint {
	return &L2TPEndpoint{
		Name:        fmt.Sprintf("%d_%s", netInfo.count, netInfo.PoolID),
		networkOpts: netInfo.networkOpts,
	}
}

func (e *L2TPEndpoint) setup() error {
	e.ConnName = e.Name + "-connection"
	e.PPPDevName = ("ppp" + e.Name)[:15]
	e.PPPOptFile = pppOptsDir + e.Name + "." + ".client"

	return nil
}

func (e *L2TPEndpoint) GetPppConfig() string {
	cfg := ""

	if e.networkOpts.PPPIPCPAcceptLocal {
		cfg += "\nipcp-accept-local"
	}
	if e.networkOpts.PPPIPCPAcceptRemote {
		cfg += "\nipcp-accept-remote"
	}
	if e.networkOpts.PPPRefuseEAP {
		cfg += "\nrefuse-eap"
	}
	if e.networkOpts.PPPRequireMSChapV2 {
		cfg += "\nrequire-mschap-v2"
	}
	if e.networkOpts.PPPNoccp {
		cfg += "\nnoccp"
	}
	if e.networkOpts.PPPNoauth {
		cfg += "\nnoauth"
	}

	cfg += fmt.Sprintf("\nifname %s", e.PPPDevName)
	cfg += fmt.Sprintf("\nname %s", e.networkOpts.PPPUsername)
	cfg += fmt.Sprintf("\npassword %s", e.networkOpts.PPPPassword)
	cfg += fmt.Sprintf("\nmtu %s", e.networkOpts.PPPMTU)
	cfg += fmt.Sprintf("\nmru %s", e.networkOpts.PPPMRU)
	cfg += fmt.Sprintf("\nidle %s", e.networkOpts.PPPIdle)
	cfg += fmt.Sprintf("\nconnect-delay %s", e.networkOpts.PPPConnectDelay)

	if e.networkOpts.PPPDebug {
		cfg += "\ndebug"
	}

	if e.networkOpts.PPPDefaultRoute {
		cfg += "\ndefaultroute"
	}
	if e.networkOpts.PPPUsepeerdns {
		cfg += "\nusepeerdns"
	}
	if e.networkOpts.PPPLock {
		cfg += "\nlock"
	}

	return cfg
}

func (e *L2TPEndpoint) GetXl2tpConfig() []string {
	return []string{
		fmt.Sprintf("lns=%s", e.networkOpts.LNSAddr),
		fmt.Sprintf("pppoptfile=%s", e.PPPOptFile),
	}
}
