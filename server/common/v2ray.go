package common

import (
	"os"

	"strconv"

	"github.com/shynome/v2ldap/v2ray"
	"github.com/thoas/go-funk"
)

// V2ray instance
var V2ray *v2ray.V2ray

func initV2ray() {
	if V2ray != nil {
		return
	}
	remoteTag, remoteGrpc, apiPortEnv, socksPortEnv := os.Getenv("RemoteTag"), os.Getenv("RemoteGrpc"), os.Getenv("V2rayAPIPort"), os.Getenv("VNEXTSocksPort")
	var apiPort, socksPort uint32
	if remoteTag == "" {
		remoteTag = "ws"
	}
	if remoteGrpc == "" {
		remoteGrpc = "127.0.0.1:3001"
	}
	if apiPortEnv == "" {
		apiPortEnv = "3001"
	}
	if port, err := strconv.Atoi(apiPortEnv); err != nil {
		panic(err)
	} else {
		apiPort = uint32(port)
	}
	if port, err := strconv.Atoi(socksPortEnv); err == nil {
		socksPort = uint32(port)
	}
	V2ray = &v2ray.V2ray{
		DB:         GetDB(),
		APIPort:    apiPort,
		RemoteTag:  remoteTag,
		RemoteGrpc: remoteGrpc,
		VNEXT:      os.Getenv("VNEXT"),
		SocksPort:  socksPort,
	}
	var users []string
	var err error
	if Ldap.BindDN != "" {
		users, err = Ldap.GetUsers()
	} else {
		dbUsers, err := V2ray.GetDBUsers()
		if err == nil {
			users = funk.Map(dbUsers, func(user v2ray.User) string { return user.Email }).([]string)
		}
	}
	if err != nil {
		panic(err)
	}
	if _, err := V2ray.Sync(users, true); err != nil {
		panic(err)
	}
}
