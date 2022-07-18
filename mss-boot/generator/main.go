package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mss-boot-io/mss-boot/core/server"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/generator/cfg"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/generator/router"
)

// @title generator API
// @version 0.0.1
// @description generator接口文档
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @host localhost:8081
// @BasePath
func main() {
	ctx := context.Background()

	r := gin.Default()
	router.Init(r.Group("/generator"))

	cfg.Cfg.Init(r)

	log.Println("starting admin manage")

	err := server.Manage.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
