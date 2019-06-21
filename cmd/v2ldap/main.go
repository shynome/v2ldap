package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/shynome/v2ldap/server"
)

var authToken = os.Getenv("token")
var authTokenHeader = "Token"

func init() {
	if authToken == "" {
		panic("env token is requried")
	}
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "7070"
	}
	addr := fmt.Sprintf("0.0.0.0:%v", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header[authTokenHeader] == nil || r.Header[authTokenHeader][0] == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`token header is requried`))
			return
		}
		token := r.Header[authTokenHeader][0]
		if token != authToken {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`token is not right`))
			return
		}
		server.APIMux.ServeHTTP(w, r)
	})

	fmt.Printf("server listen at %v \n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}

}
