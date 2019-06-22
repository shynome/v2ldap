package ldap

import (
	"net/http"

	server "github.com/shynome/v2ldap/server/common"
)

func listHandler(w http.ResponseWriter, r *http.Request) {
	users, err := server.Ldap.GetUsers()
	server.Resp(w, users, err)
}
