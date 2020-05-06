package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ping(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, resp{
		Data: "pong",
	})
}
