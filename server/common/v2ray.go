package common

import (
	"github.com/shynome/v2ldap/v2ray"
)

// V2ray instance
var V2ray *v2ray.V2ray

func initV2ray() {
	initLdap()
	V2ray = &v2ray.V2ray{
		Tag:      "",
		GrpcAddr: "",
	}
	users, err := Ldap.GetUsers()
	if err != nil {
		panic(err)
	}
	if _, err := V2ray.Sync(users, true); err != nil {
		panic(err)
	}
}
