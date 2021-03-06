/*
 * @Author: lwnmengjing
 * @Date: 2022/3/10 13:47
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/10 13:47
 */

package cfg

import (
	"net/http"
	"time"

	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/listener"
	"github.com/mss-boot-io/mss-boot/pkg/config"
	"github.com/mss-boot-io/mss-boot/pkg/config/mongodb"
	"github.com/mss-boot-io/mss-boot/pkg/config/source/s3"
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
	configSource, err := s3.New(
		s3.WithRegion("ap-northeast-1"),
		s3.WithBucket("matrix-config-center"),
		s3.WithDir("mss-boot-io/mss-boot-monorepo/admin"),
		s3.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("cfg(s3) init failed, %s\n", err.Error())
	}
	err = config.Init(configSource, &Cfg)
	//err := config.Init(Embedded, &Cfg)
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
