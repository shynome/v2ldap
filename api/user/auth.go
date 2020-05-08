package user

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func login(key []byte) func(c echo.Context) (err error) {
	return func(c echo.Context) (err error) {
		var data struct {
			SecretKey string `json:"secret_key"`
		}

		if err = c.Bind(&data); err != nil {
			return
		}

		if data.SecretKey != c.Get("token") {
			return echo.ErrUnauthorized
		}

		claims := &jwtCustomClaims{
			"admin",
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString(key)

		return c.JSON(http.StatusOK, struct {
			Token string `json:"token"`
		}{
			Token: t,
		})
	}
}

func whoami(c echo.Context) (err error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.JSON(http.StatusOK, map[string]string{
		"hello": name,
	})
}
