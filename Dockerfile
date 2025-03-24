#FROM k-harbor.buffge.com/dk/library/golang:1.24.1-alpine AS builder
FROM k-harbor.buffge.com/dk/library/golang@sha256:43c094ad24b6ac0546c62193baeb3e6e49ce14d3250845d166c77c25f64b0386 AS builder

ENV GOSUMDB=off
ENV GOPROXY=https://goproxy.buffge.com,direct
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/go/pkg/mod

ARG buildUser
ARG buildTime
ARG commitID

WORKDIR /work

COPY go.mod .
COPY go.sum .
RUN  --mount=type=cache,id=go-mod-cache,target=/go/pkg/mod,rw \
     go mod download

COPY . .
RUN  --mount=type=cache,id=go-build-cache,target=/root/.cache/go-build,rw \
     go generate ./... && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=sonic \
     -trimpath -ldflags "-s -w \
     -extldflags "-static" \
     -X github.com/real-web-world/lol-api.Commit=${commitID} \
     -X github.com/real-web-world/lol-api.BuildTime=${buildTime} \
     -X github.com/real-web-world/lol-api.BuildUser=${buildUser} \
     " -o bin/lol-api cmd/lol-api/main.go

#FROM k-harbor.buffge.com/dk/library/alpine:3.20.2
FROM k-harbor.buffge.com/dk/library/alpine@sha256:b75b7690fb4afe6fdfabfd5f1d4c8a7b710749d555bedd448dc52e9ff0dc8cc7
LABEL maintainer="buffge <admin@buffge.com>" \
      version="1.0" \
      description="hh-lol-prophet api"
ARG commitMsg
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --update --no-cache ca-certificates tzdata bash
ENV LANG=C.UTF-8 TZ=Asia/Shanghai
LABEL commitMsg=${commitMsg}

WORKDIR /work
COPY  --from=builder /work/bin /work/bin
EXPOSE 8888

ENTRYPOINT ["./bin/lol-api"]