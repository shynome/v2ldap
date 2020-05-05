package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ntoken := c.Request().Header.Get("token")
		if ntoken != c.Get("token") {
			return c.String(http.StatusUnauthorized, "token header is not right")
		}
		return next(c)
	}
}
