package v2ray

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	server "github.com/shynome/v2ldap/server/common"
	"github.com/shynome/v2ldap/v2ray"
)

var v2rayConfig []byte

func configHandler(w http.ResponseWriter, r *http.Request) {
	var config []byte
	var err error
	if v2rayConfig == nil {
		v2rayConfig, err = proto.Marshal(v2ray.Config)
	}
	config = v2rayConfig
	server.Resp(w, config, err)
}
