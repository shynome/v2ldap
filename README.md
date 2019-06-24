## 简介

这是一个将 ldap 用户同步到 v2ray 的应用

## 如何使用

运行前提: 你在一个的网络中

`curl -H 'token: yourtoken' -s 127.0.0.1:7070/v2ray/config | ddd`

## 环境变量

| name            | default value                             | desc             |
| --------------- | ----------------------------------------- | ---------------- |
| `RemoteTag`     | `ws`                                      | remtoe v2ray tag |
| `RemoteGrpc`    | `127.0.0.1:3001`                          | remote grpc addr |
| `PORT`          | `7070`                                    | http server port |
| `token`         | `uuid-long-long-long`                     | uuid             |
| `DB_PATH`       | `/app/data/v2ldap.db`                     | sqlite db path   |
| `LDAP_Host`     | `ldaps://your.company.com`                |
| `LDAP_BaseDN`   | `ou=users,dc=company,dc=com`              |
| `LDAP_Filter`   | `(&(objectclass=inetOrgPerson))`          |
| `LDAP_Attr`     | `mail`                                    |
| `LDAP_BindDN`   | `cn=v2ray-read,ou=apps,dc=company,dc=com` |
| `LDAP_Password` | `bindDNpassword`                          |

## api

| path            | params           | desc                               |
| --------------- | ---------------- | ---------------------------------- |
| `/ldap/list`    | `{}`             | 显示 ldap 服务器的用户             |
| `/v2ray/config` | `{}`             | 初始配置, 含有已存在用户的配置信息 |
| `/v2ray/sync`   | `{confirm:bool}` | 同步 ldap 的用户到 v2ray           |
| `/v2ray/uuid`   | `{email:string}` | 获取一个用户的 uuid                |
