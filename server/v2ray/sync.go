package v2ray

import (
	"net/http"

	server "github.com/shynome/v2ldap/server/common"
)

type syncHandlerParams struct {
	Confirm bool
}

func syncHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var users []string
	params := &syncHandlerParams{}
	if err = server.ParseParamsFromReq(r, &params); err != nil {
		server.Resp(w, nil, err)
		return
	}
	users, err = server.Ldap.GetUsers()
	if err != nil {
		server.Resp(w, nil, err)
		return
	}
	resp, err := server.V2ray.Sync(users, params.Confirm)
	if err != nil {
		server.Resp(w, nil, err)
		return
	}
	server.Resp(w, resp, err)
}
