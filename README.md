## 简介

v2ray 多用户管理

## 如何使用

#### 第二步: 启动服务

创建 `docker-compose.yml` 文件后, 启动服务 `docker stack deploy -c docker-compose.yml v2ray`.
如果提示集群未初始化, 使用 `docker swarm init --advertise-addr 127.0.0.1` 进行初始化

```yml
version: '3.6'

networks:
  default:
    attachable: true

services:
  ctrol:
    image: xxxxx
    volumes: ['/work/data/v2ray:/app/data']
    environment:
      PORT: 80
      token: 'uuid-long-long-long'
      VNEXTSocksPort: '1080'
    logging: &logging
      options: { max-size: '200k', max-file: '10' }
    deploy: &deploy
      endpoint_mode: dnsrr
      resources: { limits: { memory: 200M } }
      restart_policy: { condition: on-failure, max_attempts: 3 }

  client:
    image: xxxx
    command: ['v2wss']
    environment:
      APIEndpoint: 'ctrol/node?node=uuid-long-long-long'
    ports: [{ mode: host, protocol: tcp, target: 3005, published: 3005 }]
    depends_on: ['ctrol']
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

##### 启动 `v2wss`

`APIEndpoint='http://127.0.0.1:7070/node?token=your-token' v2wss`

会监听一个端口

- `3005` ws vmess port

## 环境变量

| name      | default value         | desc             |
| --------- | --------------------- | ---------------- |
| `PORT`    | `7070`                | http server port |
| `token`   | `uuid-long-long-long` | uuid             |
| `DB_PATH` | `/app/data/v2ldap.db` | sqlite db path   |

## 默认没有值的环境变量(设置了会有额外的效果)

| name             | 示例值                 | 额外效果                                |
| ---------------- | ---------------------- | --------------------------------------- |
| `VNEXTSocksPort` | `1080`                 | 如果有值则暴露一个无需认证的 socks 端口 |
| `VNEXT`          | `vmess://base64base64` | 有值的话所有流量都会使用这个节点        |
