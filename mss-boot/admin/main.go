package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mss-boot-io/mss-boot/core/server"
	"log"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/admin/cfg"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/admin/router"
)

// @title admin API
// @version 0.0.1
// @description admin接口文档
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath
func main() {
	ctx := context.Background()

	r := gin.Default()
	router.Init(r.Group("/admin"))

	cfg.Cfg.Init(r)

	log.Println("starting admin manage")

	err := server.Manage.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
