package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mss-boot-io/mss-boot/core/server"

	"github.com/WhiteMatrixTech/matrix-cloud-monorepo/oos-gateway/cfg"
	"github.com/WhiteMatrixTech/matrix-cloud-monorepo/oos-gateway/router"
)

// @title oos-gateway API
// @version 0.0.1
// @description oos-gateway接口文档
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath
func main() {
	ctx := context.Background()

	r := gin.Default()
	router.Init(r.Group("/oos-gateway"))

	cfg.Cfg.Init(r)

	log.Println("starting admin manage")

	err := server.Manage.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
