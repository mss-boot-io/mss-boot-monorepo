/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/10/24 12:54:24
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/10/24 12:54:24
 */

package handlers

import (
	"context"
	
	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task-proto/v1"
)

func (h Handler) Call(ctx context.Context, req *pb.CallRequest) (*pb.CallResponse, error) {
	result := &pb.CallResponse{
		Message: "Hello " + req.Name,
	}
	return result, nil
}
