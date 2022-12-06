package main

import (
	"context"

	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/grpc"

	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task-proto/v1"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task/cfg"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task/handlers"
)

func main() {
	ctx := context.Background()

	cfg.Cfg.Init(func(srv *grpc.Server) {
		pb.RegisterTaskServer(srv.Server(), handlers.New("task"))
	})

	log.Info("starting task manage")

	defer handlers.Cron.Stop()
	err := server.Manage.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
