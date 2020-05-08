package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shynome/v2ldap/api/node"
	"github.com/shynome/v2ldap/api/user"
)

// Register api
func Register(e *echo.Group, key []byte) {

	node.Register(e.Group("/node"))
	user.Register(e.Group("/user"), key)

}
