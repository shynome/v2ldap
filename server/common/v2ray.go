package common

import (
	"github.com/shynome/v2ldap/v2ray"
)

// V2ray instance
var V2ray *v2ray.V2ray

func initV2ray() {
	V2ray = &v2ray.V2ray{
		Tag:      "",
		GrpcAddr: "",
	}
	if _, err := V2ray.Sync(true); err != nil {
		panic(err)
	}
}
