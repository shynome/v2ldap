package v2ray

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	server "github.com/shynome/v2ldap/server/common"
)

var v2rayConfig []byte

func configHandler(w http.ResponseWriter, r *http.Request) {
	var config []byte
	var err error
	if v2rayConfig == nil {
		config := server.V2ray.GetConfig()
		v2rayConfig, err = proto.Marshal(config)
	}
	if err != nil {
		server.Resp(w, nil, err)
	}
	config = v2rayConfig
	w.Write(config)
}
