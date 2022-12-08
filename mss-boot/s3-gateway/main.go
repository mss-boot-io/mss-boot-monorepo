package main

import (
	"context"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/s3-gateway/controllers"
	"log"

	"github.com/mss-boot-io/mss-boot/core/server"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/s3-gateway/cfg"
)

// @title s3-gateway API
// @version 0.0.1
// @description s3-gateway接口文档
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath
func main() {
	ctx := context.Background()

	cfg.Cfg.Init(&controllers.OOS{})

	log.Println("starting oss-gateway manage")

	err := server.Manage.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
