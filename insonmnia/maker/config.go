package maker

import (
	"context"
	"encoding/gob"
	"io/ioutil"
	"strings"
	"time"

	"github.com/jinzhu/configor"
	log "github.com/noxiouz/zapctx/ctxlog"
	"github.com/pkg/errors"
	"github.com/sonm-io/core/accounts"
	"github.com/sonm-io/core/cmd/cli/task_config"
	"github.com/sonm-io/core/insonmnia/structs"
	pb "github.com/sonm-io/core/proto"
	"github.com/sonm-io/core/util"
	"go.uber.org/zap"
)

func init() {
	gob.Register(&task_config.YamlConfig{})
}

type Config struct {
	NumOrders    int                 `required:"true" yaml:"num_orders"`
	UpdatePeriod time.Duration       `required:"true" yaml:"update_period"`
	NodeAddr     string              `required:"true" yaml:"node_addr"`
	Eth          *accounts.EthConfig `required:"true" yaml:"ethereum"`
	DealsPath    string              `required:"true" yaml:"deals_path"`
}

// NewConfig loads a hub config from the specified YAML file.
func NewConfig(path string) (*Config, error) {
	var (
		cfg = &Config{}
		err = configor.Load(cfg, path)
	)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) GetDealFormulas() ([]*DealFormula, error) {
	ctx := context.Background()
	files, err := ioutil.ReadDir(c.DealsPath)
	if err != nil {
		return nil, err
	}

	log.G(ctx).Info("found formulas",
		zap.Int("num_files", len(files)), zap.String("path", c.DealsPath))

	idToFormula := make(map[string]*DealFormula)

	for _, file := range files {
		signature := strings.Split(file.Name(), "_")
		if len(signature) < 2 {
			log.G(ctx).Error("wrong file name", zap.String("file_name", file.Name()))
			continue
		}

		_, ok := idToFormula[signature[0]]
		if !ok {
			idToFormula[signature[0]] = &DealFormula{ID: signature[0]}
		}
		dealFormula := idToFormula[signature[0]]

		switch signature[1] {
		case "order.yaml":
			order, err := c.loadOrderFile(c.DealsPath + "/" + file.Name())
			if err != nil {
				log.G(ctx).Error("failed to load order file",
					zap.String("file_name", file.Name()),
					zap.Error(err))
				continue
			}

			dealFormula.Order = order.Unwrap()
		case "task.yaml":
			task, err := task_config.LoadConfig(c.DealsPath + "/" + file.Name())
			if err != nil {
				log.G(ctx).Error("failed to load order file",
					zap.String("file_name", file.Name()),
					zap.Error(err))
				continue
			}

			dealFormula.Task = task
		}
	}

	var out []*DealFormula
	for formulaID, dealFormula := range idToFormula {
		if dealFormula.Task == nil || dealFormula.Order == nil {
			continue
		}

		log.G(ctx).Info("Loaded formula", zap.String("formula_id", formulaID))

		out = append(out, dealFormula)
	}

	if len(out) < 1 {
		return nil, errors.New("no deal formulas found")
	}

	return out, nil
}

func (c *Config) loadOrderFile(path string) (*structs.Order, error) {
	cfg := task_config.OrderConfig{}
	err := util.LoadYamlFile(path, &cfg)
	if err != nil {
		return nil, err
	}

	order, err := cfg.IntoOrder()
	if err != nil {
		return nil, err
	}

	return order, nil
}

type DealFormula struct {
	ID    string
	Order *pb.Order
	Task  task_config.TaskConfig
}
