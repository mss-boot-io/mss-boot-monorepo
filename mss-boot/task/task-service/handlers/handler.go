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

	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task-proto/v1"
)

type Handler struct {
	handler.Handler
	pb.UnimplementedHelloworldServer
}

// New handler
func New(id string) *Handler {
	return &Handler{
		Handler: handler.Handler{
			ID:  id,
			Log: log.NewHelper(log.DefaultLogger),
		},
	}
}
