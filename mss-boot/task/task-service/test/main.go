/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/10/24 16:52:53
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/10/24 16:52:53
 */

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task-proto/v1"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %s", err.Error())
	}
	br, err := os.ReadFile("test/simple.wasm")
	if err != nil {
		log.Fatalf("could not read file: %s", err.Error())
	}
	req := &pb.AddTaskRequest{
		Spec:     "0 * * * * *",
		FuncType: pb.FuncType_FuncTypeWasm,
		FuncName: "sum",
		Content:  br,
	}

	r, err := pb.NewTaskClient(conn).AddTask(context.TODO(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Id)
}
