FROM golang:1.12.5@sha256:cf0b9f69ad1edd652a7f74a1586080b15bf6f688c545044407e28805066ef2cb as Build
WORKDIR /app
ADD go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN set -e \
  && cd cmd/v2ldap \
  && go build -o main

FROM alpine:3.9.4@sha256:769fddc7cc2f0a1c35abb2f91432e8beecf83916c421420e6a6da9f8975464b6
# 需要安装这个 ldaps 证书才可以被识别
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=Build /app/cmd//v2ldap/main /app/v2ldap
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
