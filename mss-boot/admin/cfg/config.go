/*
 * @Author: lwnmengjing
 * @Date: 2022/3/10 13:47
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/10 13:47
 */

package cfg

import (
	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/listener"
	"github.com/mss-boot-io/mss-boot/pkg/config"
	"github.com/mss-boot-io/mss-boot/pkg/config/mongodb"
	"net/http"
)

var Cfg Config

type Config struct {
	Logger   config.Logger    `yaml:"logger" json:"logger"`
	Server   config.Listen    `yaml:"server" json:"server"`
	Health   *config.Listen   `yaml:"health" json:"health"`
	Metrics  *config.Listen   `yaml:"metrics" json:"metrics"`
	Clients  config.Clients   `yaml:"clients" json:"clients"`
	Database mongodb.Database `yaml:"database" json:"database"`
	OAuth2   config.OAuth2    `yaml:"oauth2" json:"oauth2"`
}

func (e *Config) Init(handler http.Handler) {
	err := config.Init(Embedded, &Cfg)
	if err != nil {
		log.Fatalf("cfg init failed, %s\n", err.Error())
	}

	e.Logger.Init()
	e.Database.Init()

	if len(e.Clients) > 0 {
		err = e.Clients.Init()
		if err != nil {
			log.Fatalf("cfg(clients) init failed, %s\n", err.Error())
		}
	}

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
