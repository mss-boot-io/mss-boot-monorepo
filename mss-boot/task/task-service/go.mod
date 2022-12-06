module github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task

go 1.18

replace github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task-proto => ../task-proto

require (
	github.com/mss-boot-io/mss-boot-monorepo/mss-boot/task-proto v0.0.0-00010101000000-000000000000
	github.com/mss-boot-io/mss-boot/core v0.0.2
	github.com/mss-boot-io/mss-boot/pkg v0.0.3-0.20221021104239-ae78caf2490e
	github.com/robfig/cron/v3 v3.0.1
	github.com/sanity-io/litter v1.5.5
	github.com/spf13/cast v1.5.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	github.com/wasmerio/wasmer-go v1.0.4
	google.golang.org/grpc v1.51.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/go-oidc/v3 v3.2.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mss-boot-io/mss-boot/client v0.0.1 // indirect
	github.com/mss-boot-io/mss-boot/proto v0.0.1 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220407144326-9054f6ed7bac // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
