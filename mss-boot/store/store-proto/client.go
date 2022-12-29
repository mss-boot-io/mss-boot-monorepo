/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/12/13 22:03:30
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/12/13 22:03:30
 */

package store_proto

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mss-boot-io/mss-boot/core/server/grpc"
	"github.com/mss-boot-io/mss-boot/pkg/config"
	"github.com/mss-boot-io/mss-boot/pkg/config/source"

	"github.com/mss-boot-io/mss-boot-monorepo/mss-boot/store/store-proto/cfg"
	pb "github.com/mss-boot-io/mss-boot-monorepo/mss-boot/store/store-proto/v1"
)

var service = &grpc.Service{}

func init() {
	opts := make([]source.Option, 0)

	switch source.Provider(os.Getenv("config_source")) {
	case source.Local, source.FS:
		// fixme: client not support local
		opts = append(opts, source.WithProvider(source.FS),
			source.WithFrom(cfg.FS))
	case source.S3:
		_, pwd, _, _ := runtime.Caller(0)
		opts = append(opts, source.WithProvider(source.S3),
			source.WithDir(filepath.Dir(pwd)),
			source.WithProjectName("mss-boot"))
	}
	err := config.Init(&cfg.Cfg, opts...)
	if err != nil {
		log.Fatalf("cfg parse failed, %s\n", err.Error())
	}
	fmt.Println(cfg.Cfg.Client.Address, cfg.Cfg.Client.Timeout)
	err = service.Dial(cfg.Cfg.Client.Address, cfg.Cfg.Client.Timeout)
	if err != nil {
		log.Fatalf("grpc dial connect failed, %s\n", err.Error())
	}
}

// Close connection
func Close() error {
	return service.Connection.Close()
}

// GetClient get client
func GetClient() pb.StoreClient {
	return pb.NewStoreClient(service.Connection)
}
