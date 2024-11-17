.PHONY: build pack release clean
default: build

GIT_COMMIT=`git rev-list -1 HEAD`
BUILD_TIME=`TZ="Asia/Shanghai" date '+%Y-%m-%d_%H:%M:%S-%Z'`
BUILD_USER?=`whoami`
GOPROXY?=https://goproxy.buffge.com,direct
build: cmd/lol-api
	@go mod tidy && go generate ./... && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=sonic \
	-trimpath -ldflags "-s -w \
    -extldflags "-static" -X github.com/real-web-world/lol-api.Commit=${GIT_COMMIT} \
    -X github.com/real-web-world/lol-api.BuildTime=${BUILD_TIME} \
    -X github.com/real-web-world/lol-api.BuildUser=${BUILD_USER} \
    " -o bin/lol-api cmd/lol-api/main.go
doc: cmd/lol-api
	swag init -g .\cmd\lol-api\main.go
clean: bin/
	@rm -rf bin/*
upx : cmd/lol-api
	upx -9 ./bin/lol-api
