package v2ray

import (
	"github.com/thoas/go-funk"
	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/net"
	protocol "v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/inbound"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/websocket"
)

func (v2 V2ray) getPbUsers() (vmessUsers []*protocol.User) {
	var users []User
	v2.DB.Model(&User{}).Find(&users)
	vmessUsers = funk.Map(users, func(user User) *protocol.User {
		return &protocol.User{
			Email: user.Email,
			Account: serial.ToTypedMessage(&vmess.Account{
				Id:               user.UUID,
				AlterId:          64,
				SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
			}),
		}
	}).([]*protocol.User)
	return
}

// GetConfig expose v2ray config
func (v2 V2ray) GetConfig() *core.Config {

	wsPath := v2.WsPath
	if wsPath == "" {
		wsPath = "/ray"
	}
	wsPort := v2.WsPort
	if wsPort == 0 {
		wsPort = 3005
	}
	apiPort := v2.APIPort
	if apiPort == 0 {
		apiPort = 3001
	}

	if v2.config == nil {

		users := v2.getPbUsers()
		usersInbound := &core.InboundHandlerConfig{
			Tag: v2.RemoteTag,
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
				User: users,
			}),
		}

		v2.config = getV2rayConfig(apiPort)
		v2.config.Inbound = append(v2.config.Inbound, usersInbound)

	}
	return v2.config
}
