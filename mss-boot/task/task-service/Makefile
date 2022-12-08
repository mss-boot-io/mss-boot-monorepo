.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -v -ldflags "-s -w" -o application .
generate:
	go generate ./...
	#protoc --go_out=plugins=grpc:. proto/*/*.proto --go_opt=paths=source_relative
test:
	go test ./...