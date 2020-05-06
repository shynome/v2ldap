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
			Error: "数据库查询用户失败",
			Data:  err.Error(),
		})
	}
	if count == 1 {
		if err = db.Model(&model.User{}).Unscoped().Where("email = ?", params.Email).Update("deleted_at", nil).Error; err != nil {
			return c.JSON(http.StatusOK, resp{
				Error: "邮箱已存在, 进行用户恢复但失败了",
			})
		}
		return c.JSON(http.StatusOK, resp{
			Message: "邮箱已存在, 进行用户恢复成功",
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

func updateUser(c echo.Context) (err error) {
	type UpdateField struct {
		Run bool        `json:"run"`
		Val interface{} `json:"val"`
	}
	var params struct {
		ID     uint
		Update struct {
			UUID     UpdateField `json:"uuid"`
			Remark   UpdateField `json:"remark"`
			Disabled UpdateField `json:"disable"`
		} `json:"update"`
	}
	if err = c.Bind(&params); err != nil {
		return
	}
	db := model.GetDB(c)
	q := db.Where("id = ?", params.ID)
	fields := params.Update

	// 更新 UUID 值传 "0" 的话则在服务端生成随机生成 UUID
	if fields.UUID.Run {
		val, ok := fields.UUID.Val.(string)
		if ok == false {
			return c.JSON(400, resp{
				Error: "uuid 需要是 string 类型",
			})
		}
		if val == "0" {
			val = uuid.New().String()
		}
		q = q.Update("uuid = ?", val)
	}

	// 更新是否禁用
	if fields.Disabled.Run {
		val, ok := fields.Disabled.Val.(bool)
		if ok == false {
			return c.JSON(400, resp{
				Error: "disable 需要是 bool 类型",
			})
		}
		q = q.Update("disabled = ?", val)
	}

	// 更新备注
	if fields.Remark.Run {
		val, ok := fields.Remark.Val.(string)
		if ok == false {
			return c.JSON(400, resp{
				Error: "remark 需要是 string 类型",
			})
		}
		q = q.Update("remark = ?", val)
	}

	if err = q.Error; err != nil {
		return c.JSON(http.StatusOK, resp{
			Error: "更新失败",
			Data:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, resp{
		Message: "更新成功",
	})
}

func deleteUser(c echo.Context) (err error) {
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

func getUser(c echo.Context) (err error) {
	var email = c.Request().URL.Query().Get("email")
	if email == "" {
		return c.JSON(http.StatusOK, resp{
			Error: "query email field is required",
		})
	}
	db := model.GetDB(c)
	var u model.User
	if err = db.Where(&model.User{Email: email}).First(&u).Error; err != nil {
		return c.JSON(http.StatusOK, resp{
			Error: "查找用户失败",
			Data:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, resp{
		Data: u,
	})
}
