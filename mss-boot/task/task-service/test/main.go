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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task-proto/v1"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	r, err := pb.NewHelloworldClient(conn).Call(context.TODO(), &pb.CallRequest{
		Name: "lwx",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}
