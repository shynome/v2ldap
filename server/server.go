package server

import (
	"net/http"

	"github.com/shynome/v2ldap/server/common"
	"github.com/shynome/v2ldap/server/ldap"
	"github.com/shynome/v2ldap/server/v2ray"
)

// APIMux export
var APIMux = http.NewServeMux()

func init() {
	if common.Ldap.BindDN != "" {
		APIMux.Handle("/ldap/", ldap.APIMux)
	}
	APIMux.Handle("/v2ray/", v2ray.APIMux)
}
