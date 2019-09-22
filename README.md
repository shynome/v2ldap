## 简介

这是一个将 ldap 用户同步到 v2ray 的应用

`v0.3.0` 版本支持从其他数据源导入用户了, 设置 `LDAP_USERS` 即可.

## 如何使用

运行前提: 你在一个可信的网络网络环境中, 因为 v2ray grpc api port 没有做任何鉴权机制

## 自行编译

```
git clone --recursive https://github.com/shynome/v2ldap.git
cd v2ldap/cmd/v2ldap
go build -mod=vendor -o v2ldap
```

##### 启动 `v2ldap`

`source ./.env && ./v2ldap`

默认监听的端口是: `7070`
查看 ldap 用户: `curl -H 'token: yourtoken' -sSL 127.0.0.1:7070/ldap/list`

##### 启动 `v2ray`

`curl -s -H 'token: yourtoken' -sSL 127.0.0.1:7070/v2ray/config | v2ray -config=stdin: -format=pb`

会监听两个端口

- `3001` grpc api port
- `3005` ws vmess port

## 环境变量

| name            | default value                             | desc               |
| --------------- | ----------------------------------------- | ------------------ |
| `RemoteTag`     | `ws`                                      | remtoe v2ray tag   |
| `RemoteGrpc`    | `127.0.0.1:3001`                          | remote grpc addr   |
| `V2rayAPIPort`  | `3001`                                    | v2ray 远程控制端口 |
| `PORT`          | `7070`                                    | http server port   |
| `token`         | `uuid-long-long-long`                     | uuid               |
| `DB_PATH`       | `/app/data/v2ldap.db`                     | sqlite db path     |
| `LDAP_USERS`    | `file:///app/data/users.jsonc`             | 用户数据源         |
| `LDAP_Host`     | `ldaps://your.company.com`                |
| `LDAP_BaseDN`   | `ou=users,dc=company,dc=com`              |
| `LDAP_Filter`   | `(&(objectclass=inetOrgPerson))`          |
| `LDAP_Attr`     | `mail`                                    |
| `LDAP_BindDN`   | `cn=v2ray-read,ou=apps,dc=company,dc=com` |
| `LDAP_Password` | `bindDNpassword`                          |

注:

- `LDAP_USERS` 用户数据源支持 `file` 和 `https` 链接, 设置了这个就不会从 LDAP 服务器获取用户了.
  这个链接返回的内容需要是 `[]string` 的 `json`

## 默认没有值的环境变量(设置了会有额外的效果)

| name             | 示例值                 | 额外效果                                |
| ---------------- | ---------------------- | --------------------------------------- |
| `VNEXTSocksPort` | `1080`                 | 如果有值则暴露一个无需认证的 socks 端口 |
| `VNEXT`          | `vmess://base64base64` | 有值的话所有流量都会使用这个节点        |

## api

| path            | params           | desc                               |
| --------------- | ---------------- | ---------------------------------- |
| `/ldap/list`    | `{}`             | 显示 ldap 服务器的用户             |
| `/v2ray/config` | `{}`             | 初始配置, 含有已存在用户的配置信息 |
| `/v2ray/sync`   | `{confirm:bool}` | 同步 ldap 的用户到 v2ray           |
| `/v2ray/uuid`   | `{email:string}` | 获取一个用户的 uuid                |
