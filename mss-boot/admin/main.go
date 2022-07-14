package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/pkg/config"
	"log"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/admin/cfg"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/admin/router"
)

// @title admin API
// @version 0.0.1
// @description tenant接口文档
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @host localhost:9094
// @BasePath
func main() {
	c := &cfg.Config{}
	err := config.Init(flag.Lookup("c").Value.String(), c)
	if err != nil {
		log.Printf("cfg init failed, %s\n", err.Error())
		return
	}
	ctx := context.Background()

	r := gin.Default()
	router.Init(r.Group("/admin"))

	c.Init(r)

	log.Println("starting admin manage")

	err = server.Manage.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
