package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/shynome/v2ldap/api"
	"github.com/shynome/v2ldap/model"
)

var isDev = strings.HasPrefix(os.Getenv("APP_ENV"), "dev")
var authToken = os.Getenv("token")

func registerToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("token", authToken)
		return next(c)
	}
}

func main() {
	e := echo.New()

	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		if lv, err := strconv.ParseUint(lvl, 10, 8); err == nil {
			e.Logger.SetLevel(log.Lvl(lv))
		}
	}

	if authToken == "" {
		e.Logger.Fatal("env token is requried")
	}
	if err := model.Init(); err != nil {
		e.Logger.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7070"
	}

	if isDev {
		e.Use(middleware.CORS())
	}

	e.Static("/", "ui")

	api.Register(
		e.Group("/api", registerToken, model.RegisterDB),
		[]byte(authToken),
	)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))

}
