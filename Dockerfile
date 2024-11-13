FROM k-harbor.buffge.com/dk/alpine:3.20.2
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --update --no-cache ca-certificates tzdata bash
MAINTAINER buffge "admin@buffge.com"

ENV LANG=C.UTF-8 TZ=Asia/Shanghai
LABEL commitMsg=${commitMsg}

WORKDIR /work
COPY  ./bin /work/bin
EXPOSE 8888

ENTRYPOINT ["./bin/lol-api"]