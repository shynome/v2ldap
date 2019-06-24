package v2ray

import (
	"net/http"
)

// APIMux export
var APIMux = http.NewServeMux()

func init() {
	APIMux.HandleFunc("/v2ray/config", configHandler)
	APIMux.HandleFunc("/v2ray/sync", syncHandler)
}
