package main

import (
	"context"

	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/store/store-proto/v1"
	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/grpc"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/store/store-service/cfg"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/store/store-service/handlers"
)

func main() {
	ctx := context.Background()

	cfg.Cfg.Init(func(srv *grpc.Server) {
		pb.RegisterStoreServer(srv.Server(), handlers.NewStoreHandler("store"))
	})

	log.Info("starting generator manage")

	err := server.Manage.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
