/*
 * @Author: lwnmengjing
 * @Date: 2022/3/10 13:47
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/10 13:47
 */

package cfg

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/listener"
	"github.com/mss-boot-io/mss-boot/pkg/config"
	"github.com/mss-boot-io/mss-boot/pkg/config/mongodb"
	"github.com/mss-boot-io/mss-boot/pkg/config/source"
)

var (
	//go:embed *.yml
	FS  embed.FS
	Cfg Config
)

// Config 配置
type Config struct {
	Logger   config.Logger    `yaml:"logger" json:"logger"`
	Server   config.Listen    `yaml:"server" json:"server"`
	Health   *config.Listen   `yaml:"health" json:"health"`
	Metrics  *config.Listen   `yaml:"metrics" json:"metrics"`
	Database mongodb.Database `yaml:"database" json:"database"`
	OAuth2   config.OAuth2    `yaml:"oauth2" json:"oauth2"`
}

func (e *Config) Init(handler http.Handler) {
	opts := make([]source.Option, 0)
	switch source.Provider(os.Getenv("config_source")) {
	case source.FS:
		opts = append(opts, source.WithProvider(source.FS),
			source.WithFrom(FS))
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
	e.Database.Init()

	runnable := []server.Runnable{
		listener.New("admin",
			e.Server.Init(listener.WithHandler(handler))...),
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
	e.Database.Init()
	log.Info("!!! cfg change and reload")
}
