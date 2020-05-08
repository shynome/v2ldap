package user

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type resp struct {
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Register node api
func Register(g *echo.Group, key []byte) {
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: key,
		Claims:     &jwtCustomClaims{},
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			return strings.HasSuffix(path, "/login")
		},
	}))
	g.Any("/login", login(key))
	g.Any("/whoami", whoami)
	g.Any("/add", addUser)
	g.Any("/update", updateUser)
	g.Any("/delete", deleteUser)
	g.Any("/list", listUser)
	g.Any("/get", getUser)
}
