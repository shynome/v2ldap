package node

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shynome/v2ldap/model"
)

type v2rayNodeUser struct {
	Email   string
	AlterID uint64
	ID      string
}

func getUserHandler(c echo.Context) (err error) {
	utoken := c.Request().URL.Query().Get("user")
	var u model.User
	db := model.GetDB(c)
	if err = db.Where("id = ?", utoken).First(&u).Error; err != nil {
		return
	}
	var aid uint64 = 0
	if u.Version == 0 {
		aid = 64
	}
	data := v2rayNodeUser{
		Email:   u.Email,
		AlterID: aid,
		ID:      u.UUID,
	}
	return c.JSON(http.StatusOK, data)
}
