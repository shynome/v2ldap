package v2ray

import (
	"net/http"

	server "github.com/shynome/v2ldap/server/common"
	"github.com/shynome/v2ldap/v2ray"
)

type uuidHandlerParams struct {
	Email string
}
type uuidHandlerResp struct {
	UUID string `json:"uuid"`
}

func uuidHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var params uuidHandlerParams
	if err = server.ParseParamsFromReq(r, &params); err != nil {
		server.Resp(w, nil, err)
		return
	}
	var u v2ray.User
	if err = server.DB.Select([]string{"uuid"}).First(&u, "email = ?", params.Email).Error; err != nil {
		server.Resp(w, nil, err)
		return
	}
	resp := uuidHandlerResp{UUID: u.UUID}
	server.Resp(w, resp, nil)
}
