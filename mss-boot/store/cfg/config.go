package cfg

import (
	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/grpc"
	"github.com/mss-boot-io/mss-boot/core/server/listener"
	"github.com/mss-boot-io/mss-boot/pkg/config"
	"github.com/mss-boot-io/mss-boot/pkg/config/source/local"
)

var Cfg Config

// Config config
type Config struct {
	Logger   config.Logger  `yaml:"logger" json:"logger"`
	Server   config.GRPC    `yaml:"server" json:"server"`
	Health   *config.Listen `yaml:"health" json:"health"`
	Metrics  *config.Listen `yaml:"metrics" json:"metrics"`
	Cache    string         `yaml:"cache" json:"cache"`
	Queue    string         `yaml:"queue" json:"queue"`
	Locker   string         `yaml:"locker" json:"locker"`
	Provider Provider       `yaml:"provider" json:"provider"`
}

// Init init
func (e *Config) Init(handler func(srv *grpc.Server)) {
	frs, err := local.New(local.WithDir("cfg"))
	if err != nil {
		log.Fatalf("cfg init failed, %s\n", err.Error())
	}
	err = config.Init(frs, &Cfg)
	if err != nil {
		log.Fatalf("cfg init failed, %s\n", err.Error())
	}

	e.Logger.Init()
	e.Provider.Init(e.Cache, e.Queue, e.Locker)

	runnable := []server.Runnable{
		e.Server.Init(handler, grpc.WithID("store")),
	}
	if e.Health != nil {
		runnable = append(runnable, listener.NewHealthz(e.Health.Init()...))
	}
	if e.Metrics != nil {
		runnable = append(runnable, listener.NewMetrics(e.Metrics.Init()...))
	}

	server.Manage.Add(runnable...)
}

func (e *Config) OnChange() {
	e.Logger.Init()
	log.Info("!!! cfg change and reload")
}
