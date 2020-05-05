package v2ray

import (
	"os"
	"strconv"
)

// V2ray V2ray Config
type V2ray struct {
	WsPort    uint32 // ws port
	WsPath    string // ws path
	VNEXT     string // 如果有 VNEXT 将会替换默认 outbound
	SocksPort uint32 // 如果有值则暴露一个无需认证的 socks 端口
}

// New V2ray
func New() *V2ray {
	v2 := &V2ray{}
	if wsPath := os.Getenv("WsPath"); wsPath != "" {
		v2.WsPath = wsPath
	}
	if wsPort := os.Getenv("WsPort"); wsPort != "" {
		if port, err := strconv.ParseUint(wsPort, 10, 32); err == nil {
			v2.WsPort = uint32(port)
		}
	}
	if VNEXT := os.Getenv("VNEXT"); VNEXT != "" {
		v2.VNEXT = VNEXT
	}
	if socksPort := os.Getenv("VNEXTSocksPort"); socksPort != "" {
		if port, err := strconv.ParseUint(socksPort, 10, 32); err == nil {
			v2.SocksPort = uint32(port)
		}
	}
	return v2
}
