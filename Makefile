# Makefile
.PHONY: build test run clean

build:
	go build -o bin/goredis-server cmd/server/main.go
	go build -o bin/goredis-client cmd/client/main.go

test:
	go test -v ./...

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

bench:
	go test -bench=. ./...

clean:
	rm -rf bin/

docker-build:
	docker build -t goredis .

docker-run:
	docker run -p 6379:6379 goredis