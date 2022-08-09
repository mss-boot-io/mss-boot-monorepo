.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -v -ldflags "-s -w" -o application .
generate:
	go generate ./...
test:
	go test ./...