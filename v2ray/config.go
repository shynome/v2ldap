package v2ray

import (
	"github.com/shynome/v2ldap/model"
	vnext "github.com/shynome/v2rayN-vnext"
	"github.com/thoas/go-funk"
	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/net"
	protocol "v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/socks"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/inbound"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/websocket"
)

func toPbUsers(users []model.User) (vmessUsers []*protocol.User) {
	vmessUsers = funk.Map(users, func(user model.User) *protocol.User {
		return &protocol.User{
			Email: user.Email,
			Account: serial.ToTypedMessage(&vmess.Account{
				Id:               user.UUID,
				AlterId:          0,
				SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
			}),
		}
	}).([]*protocol.User)
	return
}

// GenConfig of v2ray
func (v2 V2ray) GenConfig(users []model.User) *core.Config {

	wsPath := v2.WsPath
	if wsPath == "" {
		wsPath = "/ray"
	}
	wsPort := v2.WsPort
	if wsPort == 0 {
		wsPort = 3005
	}

	pbusers := toPbUsers(users)
	usersInbound := &core.InboundHandlerConfig{
		Tag: "v2wss",
		ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
			PortRange: net.SinglePortRange(net.Port(wsPort)),
			StreamSettings: &internet.StreamConfig{
				Protocol: internet.TransportProtocol_WebSocket,
				TransportSettings: []*internet.TransportConfig{
					{
						Protocol: internet.TransportProtocol_WebSocket,
						Settings: serial.ToTypedMessage(&websocket.Config{
							Path: wsPath,
						}),
					},
				},
			},
		}),
		ProxySettings: serial.ToTypedMessage(&inbound.Config{
			User: pbusers,
		}),
	}

	config := getV2rayConfig()
	config.Inbound = append(config.Inbound, usersInbound)

	if v2.VNEXT != "" {
		if v, err := vnext.New(v2.VNEXT); err == nil {
			config.Outbound[0] = v.NewVMessOutboundConfig("direct")
		}
	}

	if v2.SocksPort != 0 {
		socksInbound := &core.InboundHandlerConfig{
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortRange: net.SinglePortRange(net.Port(v2.SocksPort)),
				Listen:    net.NewIPOrDomain(net.ParseAddress("0.0.0.0")),
			}),
			ProxySettings: serial.ToTypedMessage(&socks.ServerConfig{AuthType: socks.AuthType_NO_AUTH}),
		}
		config.Inbound = append(config.Inbound, socksInbound)
	}

	return config
}
