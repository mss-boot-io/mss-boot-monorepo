package main

import (
	"context"

	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/grpc"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task/cfg"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task/handlers"
	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task-proto/v1"
)

func main() {
	ctx := context.Background()

	cfg.Cfg.Init(func(srv *grpc.Server) {
		pb.RegisterHelloworldServer(srv.Server(), handlers.New("task"))
	})

	log.Info("starting task manage")

	err := server.Manage.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
