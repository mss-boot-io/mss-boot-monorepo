/*
 * @Author: lwnmengjing
 * @Date: 2022/3/10 13:47
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/10 13:47
 */

package cfg

import (
	"net/http"

	"github.com/mss-boot-io/mss-boot/client"
	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/listener"
	"github.com/mss-boot-io/mss-boot/pkg/config"
	"github.com/mss-boot-io/mss-boot/pkg/config/mongodb"
	"github.com/mss-boot-io/mss-boot/pkg/config/source/local"
	"github.com/mss-boot-io/mss-boot/pkg/store"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/admin/models"
)

var Cfg Config

// Config 配置
type Config struct {
	Logger   config.Logger    `yaml:"logger" json:"logger"`
	Server   config.Listen    `yaml:"server" json:"server"`
	Health   *config.Listen   `yaml:"health" json:"health"`
	Metrics  *config.Listen   `yaml:"metrics" json:"metrics"`
	Database mongodb.Database `yaml:"database" json:"database"`
	Github   Github           `yaml:"github" json:"github"`
	Clients  config.Clients   `yaml:"clients" json:"clients"`
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
	e.Database.Init()

	if len(e.Clients) > 0 {
		err = e.Clients.Init()
		if err != nil {
			log.Fatalf("cfg(clients) init failed, %s\n", err.Error())
		}
	}

	store.DefaultOAuth2Store = models.NewTenant(client.Store().GetClient())

	runnable := []server.Runnable{
		listener.New("generator",
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
