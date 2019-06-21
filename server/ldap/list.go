package ldap

import (
	"net/http"

	server "github.com/shynome/v2ldap/server/common"
	"github.com/shynome/v2ldap/v2ray"
)

func listHandler(w http.ResponseWriter, r *http.Request) {
	users, err := v2ray.Ldap.GetUsers()
	server.Resp(w, users, err)
}
