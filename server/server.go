package server

import (
	"net/http"

	"github.com/shynome/v2ldap/server/ldap"
)

// APIMux export
var APIMux = http.NewServeMux()

func init() {
	APIMux.Handle("/ldap/", ldap.APIMux)
}
