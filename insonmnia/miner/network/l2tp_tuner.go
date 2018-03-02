package network

import (
	"context"
	"net"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-plugins-helpers/ipam"
	netDriver "github.com/docker/go-plugins-helpers/network"
	log "github.com/noxiouz/zapctx/ctxlog"
	"github.com/sonm-io/core/insonmnia/structs"
	"go.uber.org/zap"
)

type L2TPTuner struct {
	cfg        *L2TPConfig
	cli        *client.Client
	netDriver  *L2TPDriver
	ipamDriver *IPAMDriver
}

func NewL2TPTuner(ctx context.Context, cfg *L2TPConfig) (*L2TPTuner, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	var (
		store = NewL2TPNetworkStore()
		tuner = &L2TPTuner{
			cfg:        cfg,
			cli:        cli,
			netDriver:  NewL2TPDriver(ctx, store),
			ipamDriver: NewIPAMDriver(ctx, store),
		}
	)

	if err := tuner.Run(ctx); err != nil {
		return nil, err
	}

	return tuner, nil
}

func (t *L2TPTuner) Run(ctx context.Context) error {
	syscall.Unlink(t.cfg.NetworkSocketPath)
	netListener, err := net.Listen("unix", t.cfg.NetworkSocketPath)
	if err != nil {
		log.G(context.Background()).Error("Failed to listen", zap.Error(err))
		return err
	}

	netHandler := netDriver.NewHandler(t.netDriver)
	if err := netHandler.Serve(netListener); err != nil {
		log.G(context.Background()).Error("Failed to serve", zap.Error(err))
	}

	syscall.Unlink(t.cfg.IPAMSocketPath)
	ipamListener, err := net.Listen("unix", t.cfg.IPAMSocketPath)
	if err != nil {
		log.G(context.Background()).Error("Failed to listen", zap.Error(err))
		return err
	}

	ipamHandler := ipam.NewHandler(t.ipamDriver)
	if err := ipamHandler.Serve(ipamListener); err != nil {
		log.G(context.Background()).Error("Failed to serve", zap.Error(err))
	}

	go func() {
		<-ctx.Done()
		log.G(context.Background()).Info("stopping tinc socket listener")
		netListener.Close()
		ipamListener.Close()
	}()

	go func() {
		log.G(ctx).Info("l2tp ipam driver has been initialized")
		ipamHandler.Serve(ipamListener)
	}()

	go func() {
		log.G(ctx).Info("l2tp network driver has been initialized")
		netHandler.Serve(netListener)
	}()

	return nil
}

func (t *L2TPTuner) Tune(net structs.Network, hostConfig *container.HostConfig, config *network.NetworkingConfig) (Cleanup, error) {
	opts := cloneOptions(net.NetworkOptions())
	createOpts := types.NetworkCreate{
		Driver:  "l2tp_net",
		Options: opts,
		IPAM:    &network.IPAM{Driver: "l2tp_ipam", Options: opts},
	}

	response, err := t.cli.NetworkCreate(context.Background(), net.ID(), createOpts)
	if err != nil {
		return nil, err
	}

	if config.EndpointsConfig == nil {
		config.EndpointsConfig = make(map[string]*network.EndpointSettings)
		config.EndpointsConfig[response.ID] = &network.EndpointSettings{
			IPAMConfig: &network.EndpointIPAMConfig{
				IPv4Address: net.NetworkAddr(),
			},
			IPAddress: net.NetworkAddr(),
			NetworkID: response.ID,
		}
	}

	return &L2TPCleaner{
		cli:       t.cli,
		networkID: response.ID,
	}, nil
}

type L2TPCleaner struct {
	networkID string
	cli       *client.Client
}

func (t *L2TPCleaner) Close() error {
	return t.cli.NetworkRemove(context.Background(), t.networkID)
}
