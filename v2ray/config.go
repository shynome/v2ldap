package v2ray

import (
	"strconv"
	"strings"

	"v2ray.com/core"
	"v2ray.com/core/app/commander"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/command"
	"v2ray.com/core/app/router"
	logLevel "v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/blackhole"
	"v2ray.com/core/proxy/dokodemo"
	"v2ray.com/core/proxy/freedom"
)

// Config v2ray first start pb config
var Config *core.Config

var (
	apiTag = "api"
)

// APIPort of v2ray grpc
var APIPort uint16 = 3001

var blackoutList = []string{
	"0.0.0.0/8",
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.0.0.0/24",
	"192.0.2.0/24",
	"192.168.0.0/16",
	"198.18.0.0/15",
	"198.51.100.0/24",
	"203.0.113.0/24",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
}

func initAppRouteBlackListRule() *router.RoutingRule {
	blackList := []*router.CIDR{}
	for _, filed := range blackoutList {

		match := strings.Split(filed, "/")
		ipStr, prefixStr := match[0], match[1]
		if ipStr == "" || prefixStr == "" {
			continue
		}
		ip := net.ParseAddress(ipStr).IP()
		prefix, err := strconv.ParseInt(prefixStr, 10, 32)
		if err != nil {
			continue
		}

		rule := &router.CIDR{
			Ip:     ip,
			Prefix: uint32(prefix),
		}
		blackList = append(blackList, rule)
	}
	return &router.RoutingRule{
		Geoip: []*router.GeoIP{
			{Cidr: blackList},
		},
		TargetTag: &router.RoutingRule_Tag{
			Tag: "blockout",
		},
	}
}

func initAppRoute() *serial.TypedMessage {

	apiRule := &router.RoutingRule{
		InboundTag: []string{apiTag},
		TargetTag:  &router.RoutingRule_Tag{Tag: apiTag},
	}

	rule := []*router.RoutingRule{apiRule}

	if len(blackoutList) != 0 {
		blackListRule := initAppRouteBlackListRule()
		rule = append(rule, blackListRule)
	}

	return serial.ToTypedMessage(&router.Config{Rule: rule})
}

func initAppDNS() *serial.TypedMessage {
	return serial.ToTypedMessage(&dns.Config{
		NameServer: []*dns.NameServer{
			{
				Address: &net.Endpoint{
					Address: &net.IPOrDomain{Address: &net.IPOrDomain_Ip{Ip: []byte{8, 8, 8, 8}}},
				},
			},
			{
				Address: &net.Endpoint{
					Address: &net.IPOrDomain{Address: &net.IPOrDomain_Ip{Ip: []byte{8, 8, 4, 4}}},
				},
			},
			{
				Address: &net.Endpoint{
					// Address: &net.IPOrDomain{Address: &net.IPOrDomain_Domain{Domain: "localhost"}},
					Address: &net.IPOrDomain{Address: &net.IPOrDomain_Ip{Ip: []byte{127, 0, 0, 1}}},
				},
			},
		},
	})
}

func initApp() []*serial.TypedMessage {
	// 开启 api 服务
	apiService := serial.ToTypedMessage(&commander.Config{
		Tag: apiTag,
		Service: []*serial.TypedMessage{
			serial.ToTypedMessage(&command.Config{}),
		},
	})
	// 设置路由
	routeService := initAppRoute()
	// dns
	dnsService := initAppDNS()
	// 设置日志
	logService := serial.ToTypedMessage(&log.Config{ErrorLogLevel: logLevel.Severity_Error})
	return []*serial.TypedMessage{
		apiService,
		dnsService,
		routeService,
		logService,
		// init
		serial.ToTypedMessage(&dispatcher.Config{}),
		serial.ToTypedMessage(&proxyman.InboundConfig{}),
		serial.ToTypedMessage(&proxyman.OutboundConfig{}),
	}
}

func initInbound() []*core.InboundHandlerConfig {
	return []*core.InboundHandlerConfig{
		// api 端口
		&core.InboundHandlerConfig{
			Tag: apiTag,
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortRange: net.SinglePortRange(net.Port(APIPort)),
				Listen:    net.NewIPOrDomain(net.LocalHostIP),
			}),
			ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
				Address:  net.NewIPOrDomain(net.LocalHostIP),
				Port:     uint32(APIPort),
				Networks: []net.Network{net.Network_TCP},
			}),
		},
	}
}

func initOutbound() []*core.OutboundHandlerConfig {
	return []*core.OutboundHandlerConfig{
		&core.OutboundHandlerConfig{
			Tag:           "direct",
			ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
		},
		&core.OutboundHandlerConfig{
			Tag: "blockout",
			ProxySettings: serial.ToTypedMessage(&blackhole.Config{
				Response: serial.ToTypedMessage(&blackhole.HTTPResponse{}),
			}),
		},
	}
}

func getV2rayConfig() *core.Config {
	app := initApp()
	inbound := initInbound()
	outbound := initOutbound()
	config := &core.Config{
		App:      app,
		Inbound:  inbound,
		Outbound: outbound,
	}
	return config
}

func initV2rayConfig() {
	Config = getV2rayConfig()
}
