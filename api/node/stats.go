package node

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func statsHandler(c echo.Context) (err error) {
	return c.String(http.StatusOK, "received")
}
