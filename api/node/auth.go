package node

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ntoken := c.Request().Header.Get("v2wss-node")
		if ntoken != c.Get("token").(string) {
			return c.String(http.StatusUnauthorized, "v2wss-node header is not right")
		}
		return next(c)
	}
}
