/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/12/2 01:22:23
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/12/2 01:22:23
 */

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cast"
	"github.com/wasmerio/wasmer-go/wasmer"

	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task-proto/v1"
)

// AddTask add task
func (h Handler) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	h.Make(ctx)

	h.Log.Debugf("AddTask: %v", req)
	var f func()
	switch req.FuncType {
	case pb.FuncType_FuncTypeWasm:
		f = func() {
			store := wasmer.NewStore(wasmer.NewEngine())
			defer store.Close()
			module, err := wasmer.NewModule(store, req.Content)
			if err != nil {
				fmt.Printf("could not create module: %s", err.Error())
				return
			}
			defer module.Close()

			instance, err := wasmer.NewInstance(module, wasmer.NewImportObject())
			if err != nil {
				fmt.Printf("could not create instance: %s", err.Error())
				return
			}
			defer instance.Close()
			nf, err := instance.Exports.GetFunction(req.FuncName)
			if err != nil {
				fmt.Printf("could not get function: %s", err.Error())
				return
			}

			fmt.Println(time.Now())
			result, err := nf()
			if req.Webhook != "" {
				u, err := url.Parse(req.Webhook)
				if err != nil {
					return
				}
				if err != nil && u != nil {
					u.Query().Set("error", err.Error())
				}
				if result != nil && u != nil {
					u.Query().Set("result", cast.ToString(result))
				}
				if u != nil {
					_, _ = http.Get(u.String())
				}
			}

		}
	default:
		return nil, fmt.Errorf("func type not support now")
	}

	h.Log.Debugf("spec: %s", req.Spec)
	id, err := Cron.AddFunc(req.Spec, f)
	if err != nil {
		h.Log.Errorf("could not add cron job: %s", err.Error())
		return nil, err
	}

	return &pb.AddTaskResponse{
		Id: int64(id),
	}, nil
}

// RemoveTask remove task
func (h Handler) RemoveTask(ctx context.Context, req *pb.RemoveTaskRequest) (*pb.RemoveTaskResponse, error) {
	h.Make(ctx)
	Cron.Remove(cron.EntryID(req.Id))
	return &pb.RemoveTaskResponse{Id: req.Id}, nil
}
