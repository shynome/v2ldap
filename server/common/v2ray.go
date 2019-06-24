package common

import (
	"os"

	"strconv"

	"github.com/shynome/v2ldap/v2ray"
)

// V2ray instance
var V2ray *v2ray.V2ray

func initV2ray() {
	if V2ray != nil {
		return
	}
	remoteTag, remoteGrpc, apiPortEnv := os.Getenv("RemoteTag"), os.Getenv("RemoteGrpc"), os.Getenv("V2rayAPIPort")
	var apiPort uint32
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
	V2ray = &v2ray.V2ray{
		DB:         GetDB(),
		APIPort:    apiPort,
		RemoteTag:  remoteTag,
		RemoteGrpc: remoteGrpc,
	}
	users, err := Ldap.GetUsers()
	if err != nil {
		panic(err)
	}
	if _, err := V2ray.Sync(users, true); err != nil {
		panic(err)
	}
}
