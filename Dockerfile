#FROM k-harbor.buffge.com/dk/library/golang:1.23.3-alpine as builder
FROM k-harbor.buffge.com/dk/library/golang@sha256:865d971152033de755444aaa5f6cb8f91c8408020e91e1457205a1f93896d9b1 as builder
ENV GOSUMDB=off
ENV GOPROXY=https://goproxy.buffge.com,direct
ARG buildUser
ARG buildTime
ARG commitID
COPY /home/runner/.cache/go-build /root/.cache/go-build
COPY . /work
WORKDIR /work
RUN  go mod tidy && go generate ./... && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=sonic \
     -trimpath -ldflags "-s -w \
     -extldflags "-static"
     -X github.com/real-web-world/lol-api.Commit=${commitID} \
     -X github.com/real-web-world/lol-api.BuildTime=${buildTime} \
     -X github.com/real-web-world/lol-api.BuildUser=${buildUser} \
     " -o bin/lol-api cmd/lol-api/main.go

#FROM k-harbor.buffge.com/dk/library/alpine:3.20.2
FROM k-harbor.buffge.com/dk/library/alpine@sha256:b75b7690fb4afe6fdfabfd5f1d4c8a7b710749d555bedd448dc52e9ff0dc8cc7
MAINTAINER buffge "admin@buffge.com"
ARG commitMsg
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --update --no-cache ca-certificates tzdata bash
ENV LANG=C.UTF-8 TZ=Asia/Shanghai
LABEL commitMsg=${commitMsg}

WORKDIR /work
COPY  --from=builder /work/bin /work/bin
EXPOSE 8888

ENTRYPOINT ["./bin/lol-api"]