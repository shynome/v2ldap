FROM golang:1.14-alpine as Build
RUN apk add --no-cache git build-base
WORKDIR /app
COPY . /app/
RUN if [ ! -f vendor/modules.txt ];then git submodule init && git submodule update; fi
RUN set -e \
  && cd cmd/v2ldap \
  && go build -mod=vendor -o main

FROM alpine
# 需要安装这个 ldaps 证书才可以被识别
RUN apk add --no-cache ca-certificates
WORKDIR /app
# 要加上这个 golang 才会处理额外添加的 dns 记录(如: --add-host)
RUN echo "hosts: files dns" > /etc/nsswitch.conf
COPY ui /app/ui
COPY --from=Build /app/cmd/v2ldap/main /app/v2ldap
VOLUME [ "/app/data" ]
ENV \
  DB_PATH='/app/data/v2ldap.db' \
  # token must be set
  token=''
CMD ["./v2ldap"]
