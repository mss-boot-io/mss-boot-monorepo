/*
 * @Author: lwnmengjing
 * @Date: 2022/3/19 1:24
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/3/19 1:24
 */

package handlers

import (
	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/store/store-proto/v1"
	log "github.com/mss-boot-io/mss-boot/core/logger"
	"github.com/mss-boot-io/mss-boot/pkg/server/handler"
)

// StoreHandler store handler
type StoreHandler struct {
	handler.Handler
	pb.UnimplementedStoreServer
}

// NewStoreHandler new store handler
func NewStoreHandler(id string) *StoreHandler {
	return &StoreHandler{
		Handler: handler.Handler{
			ID:  id,
			Log: log.NewHelper(log.DefaultLogger),
		},
	}
}
