.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -v -ldflags "-s -w" -o application .
generate:
	go generate ./...
test:
	docker run -itd --name redis-test -p 6379:6379 redis
	go test ./...