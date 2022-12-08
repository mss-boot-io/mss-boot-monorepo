/*
 * @Author: lwnmengjing
 * @Date: 2022/3/10 13:47
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/10 13:47
 */

package cfg

import (
	"net/http"

	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/listener"
	"github.com/mss-boot-io/mss-boot/pkg/config"
	"github.com/mss-boot-io/mss-boot/pkg/config/source/local"
)

var Cfg Config

// Config 配置
type Config struct {
	Logger   config.Logger  `yaml:"logger" json:"logger"`
	Server   config.Listen  `yaml:"server" json:"server"`
	Health   *config.Listen `yaml:"health" json:"health"`
	Metrics  *config.Listen `yaml:"metrics" json:"metrics"`
	Provider ProviderConfig `yaml:"provider" json:"provider"`
}

func (e *Config) Init(handler http.Handler) {
	frs, err := local.New(local.WithDir("cfg"))
	if err != nil {
		log.Fatalf("cfg init failed, %s\n", err.Error())
	}
	err = config.Init(frs, &Cfg)
	if err != nil {
		log.Fatalf("cfg init failed, %s\n", err.Error())
	}

	e.Logger.Init()
	e.Provider.Init()

	runnable := []server.Runnable{
		listener.New("s3-gateway",
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
