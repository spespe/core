package maker

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/noxiouz/zapctx/ctxlog"
	"github.com/pkg/errors"
	"github.com/sonm-io/core/cmd/cli/task_config"
	"github.com/sonm-io/core/insonmnia/node"
	pb "github.com/sonm-io/core/proto"
	"github.com/sonm-io/core/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	storagePath              = "/home/machine/maker.db"
	storageOrdersStateBucket = "orders_state"
	storageOrdersStateKey    = "state"
	storageDealsBucket       = "deals"
)

type MarketMaker struct {
	ctx                  context.Context
	cfg                  *Config
	key                  *ecdsa.PrivateKey
	ethAddr              common.Address
	dealFormulas         []*DealFormula
	ordersState          map[string]*orderInfo
	marketClient         pb.MarketClient
	taskManagementClient pb.TaskManagementClient
	dealFormulaPointer   int
	ordersStateStorage   store.Store
	dealsStorage         store.Store
}

func NewMarketMaker(ctx context.Context, cfg *Config, key *ecdsa.PrivateKey) (*MarketMaker, error) {
	_, TLSConfig, err := util.NewHitlessCertRotator(ctx, key)
	if err != nil {
		return nil, err
	}

	creds := util.NewWalletAuthenticator(util.NewTLS(TLSConfig), util.PubKeyToAddr(key.PublicKey))
	wallet, err := util.NewSelfSignedWallet(key)

	if err != nil {
		return nil, err
	}

	marketCC, err := util.MakeWalletAuthenticatedClient(ctx, creds, cfg.NodeAddr)
	if err != nil {
		return nil, err
	}

	tasksCC, err := util.MakeWalletAuthenticatedClient(ctx, creds, cfg.NodeAddr,
		grpc.WithPerRPCCredentials(util.NewWalletAccess(wallet)))
	if err != nil {
		return nil, err
	}

	dealFormulas, err := cfg.GetDealFormulas()
	if err != nil {
		return nil, err
	}

	marketMaker := &MarketMaker{
		ctx:                  ctx,
		cfg:                  cfg,
		key:                  key,
		ethAddr:              util.PubKeyToAddr(key.PublicKey),
		dealFormulas:         dealFormulas,
		ordersState:          make(map[string]*orderInfo),
		marketClient:         pb.NewMarketClient(marketCC),
		taskManagementClient: pb.NewTaskManagementClient(tasksCC),
	}

	if err := marketMaker.initStorage(); err != nil {
		return nil, err
	}

	if err := marketMaker.loadOrdersState(); err != nil {
		log.G(marketMaker.ctx).Info("starting without state", zap.String("reason", err.Error()))
	} else {
		var (
			numPendingOrders  int
			numDeployedOrders int
		)
		for _, info := range marketMaker.ordersState {
			if info.IsDeployed {
				numDeployedOrders++
			} else {
				numPendingOrders++
			}
		}
		log.G(marketMaker.ctx).Info("starting with state",
			zap.Int("num_deployed", numDeployedOrders),
			zap.Int("num_pending", numPendingOrders))
	}

	return marketMaker, nil
}

func (m *MarketMaker) Start() {

	log.G(m.ctx).Info("starting update loop", zap.Duration("update_period", m.cfg.UpdatePeriod))
	ticker := time.NewTicker(m.cfg.UpdatePeriod)

	for {
		select {
		case <-ticker.C:
			if err := m.updateOrders(); err != nil {
				log.G(m.ctx).Error("failed to update orders", zap.Error(err))
			}

			if err := m.updateTasks(); err != nil {
				log.G(m.ctx).Error("failed to update tasks", zap.Error(err))
			}

			err := m.writeOrdersState()
			if err != nil {
				log.G(m.ctx).Error("failed to dump state", zap.Error(err))
			}
		}
	}
}

func (m *MarketMaker) updateOrders() error {
	pendingOrders, err := m.getPendingOrders()
	if err != nil {
		return err
	}

	if len(pendingOrders) > 0 {
		log.G(m.ctx).Info("found pending orders", zap.Int("num_pending_orders", len(pendingOrders)))
	} else {
		log.G(m.ctx).Info("no pending orders")
	}

	if err := m.createOrders(pendingOrders); err != nil {
		return err
	}

	return nil
}

func (m *MarketMaker) getPendingOrders() ([]*pb.Order, error) {
	ordersReply, err := m.marketClient.GetOrders(m.ctx, &pb.GetOrdersRequest{
		Order: &pb.Order{
			OrderType: pb.OrderType_BID,
			ByuerID:   m.ethAddr.Hex(),
		},
		Count: uint64(m.cfg.NumOrders),
	})

	if err != nil {
		return nil, err
	}

	return ordersReply.GetOrders(), nil
}

func (m *MarketMaker) createOrders(pendingOrders []*pb.Order) error {
	if numOrdersToCreate := m.cfg.NumOrders - len(pendingOrders); numOrdersToCreate > 0 {
		log.G(m.ctx).Info("going to create orders", zap.Int("num_orders_to_create", numOrdersToCreate))

		for i := 0; i < numOrdersToCreate; i++ {
			dealFormula := m.getNextDealFormula(pendingOrders)

			order, err := m.marketClient.CreateOrder(m.ctx, dealFormula.Order)
			if err != nil {
				return err
			}

			log.G(m.ctx).Info("successfully created order",
				zap.String("order_id", order.Id),
				zap.String("formula_id", dealFormula.ID),
				zap.Int("total_orders_created", m.dealFormulaPointer))

			m.ordersState[order.Id] = &orderInfo{TaskConfig: dealFormula.Task}
		}
	}

	return nil
}

func (m *MarketMaker) updateTasks() error {
	processingReply, err := m.marketClient.GetProcessing(m.ctx, &pb.Empty{})
	if err != nil {
		return err
	}

	for _, processedOrder := range processingReply.Orders {
		if node.HandlerStatusString(uint8(processedOrder.Status)) != "Done" {
			continue
		}

		extra := strings.Split(processedOrder.GetExtra(), ": ")
		if len(extra) < 2 {
			return fmt.Errorf("failed to pasre deal ID from extra: %s", processedOrder.GetExtra())
		}
		dealID := extra[1]

		if err := m.dealsStorage.Put(time.Now().String(), []byte(dealID), nil); err != nil {
			log.G(m.ctx).Error("failed to write deal to disk",
				zap.Error(err))
		}

		orderInfo, ok := m.ordersState[processedOrder.Id]
		if !ok {
			return fmt.Errorf("failed to get task config for order %s, deal %s", processedOrder.Id, dealID)
		}

		if orderInfo.IsDeployed {
			return nil
		}

		log.G(m.ctx).Info("order matched, trying to start task",
			zap.String("order_id", processedOrder.Id),
			zap.String("deal_id", dealID))

		if err := m.startTask(dealID, orderInfo.TaskConfig); err != nil {
			return err
		}

		orderInfo.IsDeployed = true
	}

	return nil
}

func (m *MarketMaker) startTask(dealID string, taskConfig task_config.TaskConfig) error {
	deal := &pb.Deal{
		Id:      dealID,
		BuyerID: util.PubKeyToAddr(m.key.PublicKey).Hex(),
	}

	taskStartReply, err := m.taskManagementClient.Start(m.ctx, &pb.HubStartTaskRequest{
		Deal:          deal,
		Image:         taskConfig.GetImageName(),
		Registry:      taskConfig.GetRegistryName(),
		Auth:          taskConfig.GetRegistryAuth(),
		PublicKeyData: taskConfig.GetSSHKey(),
		Env:           taskConfig.GetEnvVars(),
	})
	if err != nil {
		return errors.Wrap(err, "failed to start a task")
	}

	log.G(m.ctx).Info("started a task",
		zap.String("task_id", taskStartReply.Id),
		zap.Any("endpoint", taskStartReply.Endpoint),
		zap.String("hub_addr", taskStartReply.HubAddr))

	return nil
}

func (m *MarketMaker) getNextDealFormula(activeOrders []*pb.Order) *DealFormula {
	m.dealFormulaPointer++

	return m.dealFormulas[m.dealFormulaPointer%len(m.dealFormulas)]
}

func (m *MarketMaker) writeOrdersState() error {
	var (
		buf bytes.Buffer
		enc = gob.NewEncoder(&buf)
		err = enc.Encode(m.ordersState)
	)
	if err != nil {
		return errors.Wrap(err, "failed to encode")
	}

	if err := m.ordersStateStorage.Put(storageOrdersStateKey, buf.Bytes(), nil); err != nil {
		return err
	}

	return nil
}

func (m *MarketMaker) loadOrdersState() error {
	kvPair, err := m.ordersStateStorage.Get(storageOrdersStateKey)
	if err != nil {
		return err
	}

	var (
		dst = make(map[string]*orderInfo)
		rdr = bytes.NewReader(kvPair.Value)
		dec = gob.NewDecoder(rdr)
	)
	if err = dec.Decode(&dst); err != nil {
		return err
	}

	m.ordersState = dst

	return nil
}

func (m *MarketMaker) initStorage() error {
	boltdb.Register()

	log.G(m.ctx).Info("creating storage", zap.Any("storage_path", storagePath))

	endpoints := []string{storagePath}
	backend := store.Backend(store.BOLTDB)

	ordersStateStorage, err := libkv.NewStore(backend, endpoints, &store.Config{
		TLS:    nil,
		Bucket: storageOrdersStateBucket,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create ordersState storage")
	}

	m.ordersStateStorage = ordersStateStorage

	dealsStorage, err := libkv.NewStore(backend, endpoints, &store.Config{
		TLS:    nil,
		Bucket: storageDealsBucket,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create deals storage")
	}

	m.dealsStorage = dealsStorage

	return nil
}

type orderInfo struct {
	IsDeployed bool
	TaskConfig task_config.TaskConfig
}
