package node

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// Register node api
func Register(g *echo.Group) {
	g.Use(auth)
	g.Any("/", func(c echo.Context) (err error) {
		action := c.Request().Header.Get("v2wss-action")
		switch action {
		case "GetPbConfig":
			return configHandler(c)
		case "GetUser":
			return getUserHandler(c)
		case "PushStats":
			return statsHandler(c)
		default:
			return fmt.Errorf("no action for you")
		}
	})
}
