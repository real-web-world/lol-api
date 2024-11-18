FROM k-harbor.buffge.com/dk/library/golang:1.23.3-alpine as builder
ENV GOPROXY=https://goproxy.buffge.com,direct
ARG buildUser
ARG buildTime
ARG commitID
ARG commitMsg

COPY . /work
WORKDIR /work
RUN  go mod tidy && go generate ./... && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=sonic \
     -trimpath -ldflags "-s -w \
     -extldflags "-static" -X github.com/real-web-world/lol-api.Commit=${commitID} \
     -X github.com/real-web-world/lol-api.BuildTime=${buildTime} \
     -X github.com/real-web-world/lol-api.BuildUser=${buildUser} \
     " -o bin/lol-api cmd/lol-api/main.go

FROM k-harbor.buffge.com/dk/library/alpine:3.20.2
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --update --no-cache ca-certificates tzdata bash
MAINTAINER buffge "admin@buffge.com"

ENV LANG=C.UTF-8 TZ=Asia/Shanghai
LABEL commitMsg=${commitMsg}

WORKDIR /work
COPY  --from=builder /work/bin /work/bin
EXPOSE 8888

ENTRYPOINT ["./bin/lol-api"]