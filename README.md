## 简介

这是一个将 ldap 用户同步到 v2ray 的应用

`v0.3.0` 版本支持从其他数据源导入用户了, 设置 `LDAP_USERS` 即可.

注: 用户数量不能太多, 因为同步用户的时候是直接将所有用户加载到程序里进行对比哪些用户要增减的, 估计 20 以下是没问题的

## 如何使用

运行前提: 你在一个可信的网络网络环境中, 因为 v2ray grpc api port 没有做任何鉴权机制

#### 第一步: 创建用户

首先我们假定保存数据的文件夹为 `/work/data/v2ray` 用 `${V2rayRoot}` 代称. 若不存在自行创建 或者 使用其他的文件夹也是可以的

创建用户定义文件: `${V2rayRoot}/users.jsonc`

```jsonc
[
  "v2ray@v2ray.com", // v2ray 的第一个用户
  "v2ray2@v2ray.com" // 第二个用户. 虽然任意文本都可以但是建议使用邮箱, 因为后续可能会以此为基础添加一些额外的功能
]
```

#### 第二步: 启动服务

创建 `docker-compose.yml` 文件后, 启动服务 `docker stack deploy -c docker-compose.yml v2ray`.
如果提示集群未初始化, 使用 `docker swarm init --advertise-addr 127.0.0.1` 进行初始化

```yml
version: "3.6"

networks:
  default:
    attachable: true

services:
  ctrol:
    image: shynome/v2ldap:0.3.2@sha256:e761daac2973ae679a70d19b315ac5830871a2df778802a7409c14970ebc7f7d
    volumes: ["/work/data/v2ray:/app/data"]
    environment:
      PORT: 80
      token: "uuid-long-long-long"
      RemoteGrpc: "client:3001"
      VNEXTSocksPort: "1080"
    logging: &logging
      options: { max-size: "200k", max-file: "10" }
    deploy: &deploy
      endpoint_mode: dnsrr
      resources: { limits: { memory: 200M } }
      restart_policy: { condition: on-failure, max_attempts: 3 }

  client:
    image: shynome/v2ray-only:4.19.1@sha256:160f856ff553da8b043ea23df3c8937589cf9ac6a0e2d76e7c0093bf3aa5d6ad
    command:["-c",'wget -q -O - --header="token: uuid-long-long-long" http://ctrol/v2ray/config | v2ray -format=pb -config=stdin:']
    ports: [{ mode: host, protocol: tcp, target: 3005, published: 3005 }]
    depends_on: ["ctrol"]
    logging: { <<: *logging }
    deploy: { <<: *deploy }
```

检查 `client` 是否启动成功

```sh
docker run --rm -ti --net v2ray_default byrnedo/alpine-curl -x socks5://client:1080 google.com
```

使用 `Caddy` 进行 `ws` 转发, 配置 `tls`+`ws`

```conf
# Caddy 示例配置文件
domain.com {

  # 目前 ws 固定路径是 `/ray` , 后续版本会允许进行调整. -(不知道什么时候有空改,有 pr 就好了 x_x)-
  proxy /ray 127.0.0.1:3005 {
    websocket
    header_upstream -Origin
  }

}
```

#### 第四步: 控制

```sh
alias curl='docker run --rm -ti --net v2ray_default byrnedo/alpine-curl -H "token: uuid-long-long-long"'
# 查看用户列表
curl -d '{}' ctrol/ldap/list
# 获取用户的 uuid
curl -d '{"email":"v2ray@v2ray.com"}' ctrol/v2ray/uuid
# 查看进行用户同步后会删减的用户. 用户同步是和 `${V2rayRoot}/users.jsonc` 这个文件中定义的用户进行同步
curl -d '{}' ctrol/v2ray/sync
# 进行用户同步
curl -d '{"confirm":true}' ctrol/v2ray/sync
```

#### 第三步: v2ray 客户端配置

| name         | value                   |
| ------------ | ----------------------- |
| 服务器       | domain.com              |
| 端口         | 443                     |
| 额外 ID      | 上一步获取到的用户 uuid |
| 加密协议     | auto                    |
| 传输协议     | ws                      |
| 别名         | domain com              |
| 伪装类型     | none                    |
| 伪装域名     | (留空)                  |
| 路径         | /ray                    |
| 底层传输安全 | tls                     |

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
| `LDAP_USERS`    | `file:///app/data/users.jsonc`            | 用户数据源         |
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
