/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/10/24 12:39:22
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/10/24 12:39:22
 */

package handlers

import (
	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/pkg/server/handler"
	"github.com/robfig/cron/v3"

	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task-proto/v1"
)

var Cron = cron.New(cron.WithSeconds())

type Handler struct {
	handler.Handler
	pb.UnimplementedTaskServer
}

// New handler
func New(id string) *Handler {
	Cron.Start()
	return &Handler{
		Handler: handler.Handler{
			ID:  id,
			Log: log.NewHelper(log.DefaultLogger),
		},
	}
}
