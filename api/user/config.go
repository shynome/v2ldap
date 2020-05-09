package user

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func config(c echo.Context) (err error) {
	d := struct {
		WsURL string `json:"ws_url"`
	}{
		WsURL: os.Getenv("WsURL"),
	}
	return c.JSON(http.StatusOK, resp{
		Data: d,
	})
}
