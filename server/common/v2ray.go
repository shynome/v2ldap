package common

import (
	"github.com/shynome/v2ldap/v2ray"
)

// V2ray instance
var V2ray *v2ray.V2ray

func initV2ray() {
	if V2ray.DB != nil {
		return
	}
	V2ray = &v2ray.V2ray{
		DB: GetDB(),
	}
	users, err := Ldap.GetUsers()
	if err != nil {
		panic(err)
	}
	if _, err := V2ray.Sync(users, true); err != nil {
		panic(err)
	}
}
