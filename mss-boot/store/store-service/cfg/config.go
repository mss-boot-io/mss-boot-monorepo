package cfg

import (
	"os"
	"path/filepath"
	"runtime"

	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/grpc"
	"github.com/mss-boot-io/mss-boot/core/server/listener"
	"github.com/mss-boot-io/mss-boot/pkg/config"
	"github.com/mss-boot-io/mss-boot/pkg/config/source"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/store/store-proto/cfg"
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
	opts := make([]source.Option, 0)
	switch source.Provider(os.Getenv("config_source")) {
	case source.FS:
		opts = append(opts, source.WithProvider(source.FS),
			source.WithFrom(cfg.FS))
	case source.Local:
		opts = append(opts, source.WithProvider(source.Local),
			source.WithDir("cfg"))
	case source.S3:
		_, pwd, _, _ := runtime.Caller(1)
		opts = append(opts, source.WithProvider(source.S3),
			source.WithDir(filepath.Dir(pwd)),
			source.WithProjectName("mss-boot"))
	}
	err := config.Init(e, opts...)
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
