package main

import (
	"context"
	"flag"
	"os"

	log "github.com/noxiouz/zapctx/ctxlog"
	"github.com/sonm-io/core/insonmnia/logging"
	"github.com/sonm-io/core/insonmnia/maker"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "maker.yaml", "Path to market maker config file")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	cfg, err := maker.NewConfig(*configPath)
	if err != nil {
		log.GetLogger(ctx).Error("failed to load config", zap.Error(err))
		os.Exit(1)
	}

	key, err := cfg.Eth.LoadKey()
	if err != nil {
		log.GetLogger(ctx).Error("failed load private key", zap.Error(err))
		os.Exit(1)
	}

	logger := logging.BuildLogger(-1, true)
	ctx = log.WithLogger(context.Background(), logger)

	mk, err := maker.NewMarketMaker(ctx, cfg, key)
	if err != nil {
		log.G(ctx).Error("cannot start Locator service", zap.Error(err))
		os.Exit(1)
	}

	mk.Start()
}
