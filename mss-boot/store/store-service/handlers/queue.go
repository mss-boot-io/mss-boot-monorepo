/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/12/15 16:21:43
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/12/15 16:21:43
 */

package handlers

import (
	"context"
	"encoding/json"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/store/store-service/pkg/storage/queue"

	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/store/store-proto/v1"
	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/store/store-service/pkg/storage"
)

func (e StoreHandler) Append(c context.Context, req *pb.AppendReq) (*pb.AppendResp, error) {
	e.Make(c)
	e.Log.Debugf("Append: %v", req)
	resp := &pb.AppendResp{}
	message := &queue.Message{}
	message.SetPrefix(req.Prefix)
	message.SetStream(req.Stream)
	var values map[string]any
	err := json.Unmarshal(req.Values, &values)
	if err != nil {
		e.Log.Errorf("json.Unmarshal: %v", err)
		return resp, err
	}
	message.SetValues(values)
	message.SetID(req.Id)
	err = storage.Queue.Append(message)
	if err != nil {
		e.Log.Errorf("storage.Queue.Append: %v", err)
	}
	return nil, nil
}

func (e StoreHandler) Register(req *pb.RegisterReq, stream pb.Store_RegisterServer) error {
	e.Log.Debugf("Register: %v", req)

	storage.Queue.Register(req.Stream, func(msg storage.Messager) error {
		values := msg.GetValues()
		rb, err := json.Marshal(values)
		if err != nil {
			e.Log.Errorf("json.Marshal: %v", err)
			return err
		}
		return stream.Send(&pb.RegisterResp{
			Id:     msg.GetStream(),
			Stream: msg.GetStream(),
			Prefix: msg.GetPrefix(),
			Values: rb,
		})
	})
	<-stream.Context().Done()
	return nil
}
