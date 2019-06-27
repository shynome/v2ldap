FROM golang:1.12.5-alpine@sha256:bf3243ef1ddd18d190f22ab58c08750a3ded13c0d06a6a2a6f7e4c3451177dc4 as Build
RUN apk add --no-cache git build-base
WORKDIR /app
COPY . /app/
RUN if [ ! -f vendor/modules.txt ];then git submodule init && git submodule update; fi
RUN set -e \
  && cd cmd/v2ldap \
  && go build -mod=vendor -o main

FROM alpine:3.9.4@sha256:769fddc7cc2f0a1c35abb2f91432e8beecf83916c421420e6a6da9f8975464b6
# 需要安装这个 ldaps 证书才可以被识别
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=Build /app/cmd//v2ldap/main /app/v2ldap
# 要加上这个 golang 才会处理额外添加的 dns 记录(如: --add-host)
RUN echo "hosts: files dns" > /etc/nsswitch.conf
VOLUME [ "/app/data" ]
ENV \
  DB_PATH='/app/data/v2ldap.db' \
  # token must be set
  token='' \
  # LDAP must be set
  LDAP_Host='ldaps://your.comany.com' \
  LDAP_BaseDN='ou=users,dc=fevergroup,dc=com' \
  LDAP_Filter='(&(objectclass=inetOrgPerson))' \
  LDAP_Attr='mail' \
  LDAP_BindDN='' \
  LDAP_Password=''
CMD ["./v2ldap"]
