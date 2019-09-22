0.3.2
------------
#### CHANGE
- `LDAP_USERS` 现在支持 `jsonc` 格式的内容了

#### BREAK
- `LDAP_USERS` 现在默认值是 `file:///app/data/users.jsonc` 了, 如果你使用了 `0.3.1` 版本的话需要手动覆盖这个默认值

0.3.1
------------
#### FIX
- `LDAP_USERS` 忘记在 `Dockerfile` 中设置了, 现在加上

0.3.0
------------
#### ADD
- 添加了环境变量 `LDAP_USERS` 用户数据源, 支持 `file` 和 `https` 链接, 设置了这个就不会从 `LDAP` 服务器获取用户了.
  这个链接返回的内容需要是 `[]string` 的 `json`

0.2.3
------------
#### 优化
- `go build -mod=vendor` 上传 vendor 文件夹固化依赖, 保证每个人都能 build 成功

0.2.2
------------
#### fix
- 修复 `docker run --add-host` 添加的 `host` 不被 `golang` 程序处理的问题

0.2.1
------------
#### fix
- `VNEXTSocksPort` 现在不设置也不会报错了

0.2.0
------------
#### add
- `v2ray` 添加了 `SocksPort` 设置, 如果设置了该值的话, `GetConfig` 的配置将会含有一个无需认证的 `socks inbound`

0.1.0
------------
#### add
- 添加 `VNEXT` 支持

