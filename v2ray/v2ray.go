package v2ray

// V2ray V2ray Config
type V2ray struct {
	WsPort    uint32 // ws port
	WsPath    string // ws path
	VNEXT     string // 如果有 VNEXT 将会替换默认 outbound
	SocksPort uint32 // 如果有值则暴露一个无需认证的 socks 端口
}
