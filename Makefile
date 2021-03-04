pwd=$(shell pwd)

all: server 

server:
	PROJ_DIR=${pwd} GOFLAGS=-mod=vendor go run ./main.go server

docker.build:
	docker build -f build/docker/Dockerfile -t httpserver .

go-vendor:
	go mod tidy
	go mod vendor
