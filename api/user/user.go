package user

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/shynome/v2ldap/model"
)

func addUser(c echo.Context) (err error) {
	var params struct {
		Email string
	}
	if err = c.Bind(&params); err != nil {
		return
	}
	db := model.GetDB(c)

	var count int64
	if err = db.Model(&model.User{}).Unscoped().Where("email = ?", params.Email).Count(&count).Error; err != nil {
		return c.JSON(http.StatusOK, resp{
			Error: "数据库查询失败",
			Data:  err.Error(),
		})
	}
	if count == 1 {
		if err = db.Model(&model.User{}).Unscoped().Where("email = ?", params.Email).Update("deleted_at", nil).Error; err != nil {
			return c.JSON(http.StatusOK, resp{
				Error: "激活用户失败",
			})
		}
		return c.JSON(http.StatusOK, resp{
			Message: "已激活用户",
		})
	}

	var u = model.User{
		UUID:    uuid.New().String(),
		Email:   params.Email,
		Version: 1,
	}
	if err = db.Create(&u).Error; err != nil {
		return c.JSON(http.StatusOK, resp{
			Error: "添加失败",
			Data:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, resp{
		Message: "添加成功",
	})
}

func disableUser(c echo.Context) (err error) {
	var params struct {
		Email string
	}
	if err = c.Bind(&params); err != nil {
		return
	}
	db := model.GetDB(c)
	if err = db.Where("email = ?", params.Email).Delete(&model.User{}).Error; err != nil {
		return c.JSON(http.StatusOK, resp{
			Error: "删除失败",
			Data:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, resp{
		Message: "删除成功",
	})
}

func listUser(c echo.Context) (err error) {
	var options struct {
		Disabled bool
		All      bool
	}
	if err = c.Bind(&options); err != nil {
		return
	}
	db := model.GetDB(c)
	var users []model.User
	var q = db
	if options.All {
		q = q.Unscoped()
	} else if options.Disabled {
		q = q.Unscoped().Not("deleted_at", gorm.Expr("NULL"))
	}
	if err = q.Find(&users).Error; err != nil {
		return c.JSON(http.StatusOK, resp{
			Message: "查询用户失败",
			Data:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, resp{
		Data: users,
	})
}
