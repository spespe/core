package network

import (
	"context"
	"io/ioutil"
	"os/exec"
	"time"

	"os"

	"net"

	"github.com/docker/go-plugins-helpers/ipam"
	log "github.com/noxiouz/zapctx/ctxlog"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type IPAMDriver struct {
	ctx     context.Context
	counter int
	store   *L2TPNetworkStore
}

func NewIPAMDriver(ctx context.Context, store *L2TPNetworkStore) *IPAMDriver {
	return &IPAMDriver{
		ctx:   ctx,
		store: store,
	}
}

func (d *IPAMDriver) RequestPool(request *ipam.RequestPoolRequest) (*ipam.RequestPoolResponse, error) {
	log.G(d.ctx).Info("received RequestPool request", zap.Any("request", request))

	opts, err := parseOptsIPAM(request)
	if err != nil {
		log.G(d.ctx).Error("failed to parse options", zap.Error(err))
		return nil, errors.Wrap(err, "failed to parse options")
	}

	netInfo := newNetworkInfo(opts)
	if err := netInfo.Setup(); err != nil {
		log.G(d.ctx).Error("failed to setup network", zap.Error(err))
		return nil, err
	}

	if err := d.store.AddNetwork(netInfo.PoolID, netInfo); err != nil {
		log.G(d.ctx).Error("failed to add network to store", zap.Error(err))
		return nil, err
	}

	return &ipam.RequestPoolResponse{PoolID: netInfo.PoolID, Pool: opts.Subnet}, nil
}

func (d *IPAMDriver) RequestAddress(request *ipam.RequestAddressRequest) (*ipam.RequestAddressResponse, error) {
	log.G(d.ctx).Info("received RequestAddress request", zap.Any("request", request))

	if len(request.Address) > 0 {
		log.G(d.ctx).Error("requests for specific addresses are not supported")
		return nil, errors.New("requests for specific addresses are not supported")
	}

	netInfo, err := d.store.GetNetwork(request.PoolID)
	if err != nil {
		log.G(d.ctx).Error("failed to get network", zap.String("pool_id", request.PoolID), zap.Error(err))
		return nil, errors.Wrap(err, "failed to get network")
	}

	eptInfo := NewL2TPEndpoint(netInfo)
	if err := eptInfo.setup(); err != nil {
		log.G(d.ctx).Error("failed to setup endpoint", zap.String("pool_id", netInfo.PoolID),
			zap.String("network_id", netInfo.ID), zap.Error(err))
		return nil, err
	}

	netInfo.ConnInc()

	var (
		pppCfg       = eptInfo.GetPppConfig()
		xl2tpdCfg    = eptInfo.GetXl2tpConfig()
		addCfgCmd    = exec.Command("xl2tpd-control", "add", eptInfo.ConnName, xl2tpdCfg[0], xl2tpdCfg[1])
		setupConnCmd = exec.Command("xl2tpd-control", "connect", eptInfo.ConnName)
	)
	log.G(d.ctx).Info("creating ppp options file", zap.String("network_id", netInfo.ID),
		zap.String("ppo_opt_file", eptInfo.PPPOptFile))
	if err := ioutil.WriteFile(eptInfo.PPPOptFile, []byte(pppCfg), 0644); err != nil {
		log.G(d.ctx).Error("failed to create ppp options file", zap.String("network_id", netInfo.ID),
			zap.Any("config", xl2tpdCfg), zap.Error(err))
		return nil, errors.Wrapf(err, "failed to create ppp options file for network %s, config is `%s`",
			netInfo.ID, xl2tpdCfg)
	}

	log.G(d.ctx).Info("adding xl2tp connection config", zap.String("network_id", netInfo.ID),
		zap.String("endpoint_name", eptInfo.Name), zap.Any("config", xl2tpdCfg))
	if err := addCfgCmd.Run(); err != nil {
		log.G(d.ctx).Error("failed to add xl2tpd config", zap.String("network_id", netInfo.ID),
			zap.Any("config", xl2tpdCfg), zap.Error(err))
		return nil, errors.Wrapf(err, "failed to add xl2tpd connection config for network %s, config is `%s`",
			netInfo.ID, xl2tpdCfg)
	}

	log.G(d.ctx).Info("setting up xl2tpd connection", zap.String("connection_name", eptInfo.ConnName),
		zap.String("network_id", netInfo.ID), zap.String("endpoint_name", eptInfo.Name))
	if err := setupConnCmd.Run(); err != nil {
		log.G(d.ctx).Error("xl2tpd failed to setup connection", zap.String("network_id", netInfo.ID),
			zap.Any("config", xl2tpdCfg), zap.Error(err))
		return nil, errors.Wrapf(err, "failed to add xl2tpd config for network %s, config is `%s`",
			netInfo.ID, xl2tpdCfg)
	}

	assignedCIDR, err := d.getAssignedCIDR(eptInfo.PPPDevName)
	if err != nil {
		log.G(d.ctx).Error("failed to get assigned IP", zap.String("network_id", netInfo.ID),
			zap.Any("config", xl2tpdCfg), zap.Error(err))
		return nil, errors.Wrap(err, "failed to get assigned IP")
	}

	log.G(d.ctx).Info("received IP", zap.String("network_id", netInfo.ID),
		zap.String("ip", assignedCIDR))

	eptInfo.AssignedIP, _ = getAddrFromCIDR(assignedCIDR)
	eptInfo.AssignedCIDR = assignedCIDR

	if err := netInfo.store.AddEndpoint(eptInfo.AssignedIP, eptInfo); err != nil {
		log.G(d.ctx).Error("failed to add endpoint", zap.Error(err))
		return nil, err
	}

	return &ipam.RequestAddressResponse{Address: eptInfo.AssignedCIDR}, nil
}

func (d *IPAMDriver) ReleasePool(request *ipam.ReleasePoolRequest) error {
	log.G(d.ctx).Info("received ReleasePool request", zap.Any("request", request))

	netInfo, err := d.store.GetNetwork(request.PoolID)
	if err != nil {
		log.G(d.ctx).Error("failed to get network", zap.String("pool_id", request.PoolID), zap.Error(err))
		return errors.Wrap(err, "failed to get network info")
	}

	// Normally there won't be any endpoints: `ReleasePool` is called after releasing all addresses.
	for _, eptInfo := range netInfo.store.GetEndpoints() {
		if err := d.removeEndpoint(netInfo, eptInfo); err != nil {
			log.G(d.ctx).Error("xl2tpd failed to removeEndpoint", zap.String("pool_id", netInfo.PoolID),
				zap.String("network_id", netInfo.ID), zap.Error(err))
		}
	}

	if err := d.store.RemoveNetwork(request.PoolID); err != nil {
		log.G(d.ctx).Error("failed to remove network", zap.String("pool_id", netInfo.PoolID),
			zap.String("network_id", netInfo.ID), zap.Error(err))
		return err
	}

	return nil
}

func (d *IPAMDriver) ReleaseAddress(request *ipam.ReleaseAddressRequest) error {
	log.G(d.ctx).Info("received ReleaseAddress request", zap.Any("request", request))

	netInfo, err := d.store.GetNetwork(request.PoolID)
	if err != nil {
		log.G(d.ctx).Error("failed to get network", zap.String("pool_id", request.PoolID), zap.Error(err))
		return errors.Wrap(err, "failed to get network")
	}

	eptInfo, err := netInfo.store.GetEndpoint(request.Address)
	if err != nil {
		log.G(d.ctx).Error("failed to get endpoint", zap.String("pool_id", request.PoolID),
			zap.String("network_id", netInfo.ID), zap.String("ip", request.Address), zap.Error(err))
		return errors.Wrap(err, "failed to get endpoint")
	}

	if err := d.removeEndpoint(netInfo, eptInfo); err != nil {
		log.G(d.ctx).Error("xl2tpd failed to removeEndpoint", zap.String("pool_id", netInfo.PoolID),
			zap.String("network_id", netInfo.ID), zap.Error(err))
		return errors.Wrap(err, "xl2tpd failed to removeEndpoint")
	}

	return nil
}

func (d *IPAMDriver) GetCapabilities() (*ipam.CapabilitiesResponse, error) {
	log.G(d.ctx).Info("received GetCapabilities request")
	return &ipam.CapabilitiesResponse{RequiresMACAddress: false}, nil
}

func (d *IPAMDriver) GetDefaultAddressSpaces() (*ipam.AddressSpacesResponse, error) {
	log.G(d.ctx).Info("received GetDefaultAddressSpaces request")
	return &ipam.AddressSpacesResponse{}, nil
}

func (d *IPAMDriver) getAssignedCIDR(devName string) (string, error) {
	time.Sleep(time.Second * 5)
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range ifaces {
		if i.Name == devName {
			addrs, err := i.Addrs()
			if err != nil {
				return "", err
			}

			if len(addrs) < 1 {
				return "", errors.New("no addresses assigned!")
			}

			return addrs[0].String(), nil
		}
	}

	return "", errors.Errorf("device %s not found", devName)
}

func (d *IPAMDriver) removeEndpoint(netInfo *L2TPNetwork, eptInfo *L2TPEndpoint) error {
	disconnectCmd := exec.Command("xl2tpd-control", "disconnect", eptInfo.ConnName)
	if err := disconnectCmd.Run(); err != nil {
		return errors.Wrapf(err, "xl2rpd failed to close connection %s", eptInfo.ConnName)
	}

	if err := os.Remove(eptInfo.PPPOptFile); err != nil {
		return errors.Wrapf(err, "failed to remove ppp opts file %s", eptInfo.PPPOptFile)
	}

	if err := netInfo.store.RemoveEndpoint(eptInfo.AssignedIP); err != nil {
		return errors.Wrap(err, "failed to remove endpoint from store")
	}

	return nil
}
