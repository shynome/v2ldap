package node

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/shynome/v2ldap/model"
	"github.com/shynome/v2ldap/v2ray"
)

func configHandler(c echo.Context) (err error) {
	db := model.GetDB(c)
	v2 := &v2ray.V2ray{}
	var users = []model.User{}
	if err = db.Unscoped().Where("version = ?", gorm.Expr("NULL")).Find(&users).Error; err != nil {
		return c.String(500, "查找用户出错")
	}
	config := v2.GenConfig(users)
	pbconfig, err := proto.Marshal(config)
	if err != nil {
		return
	}
	return c.Blob(http.StatusOK, "application/v2ray-pbconfig", pbconfig)
}
