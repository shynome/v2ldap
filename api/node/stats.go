package node

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/shynome/v2ldap/model"
)

func saveStats(db *gorm.DB, logger echo.Logger, stats [][]interface{}) {
	var saveStat = func(stat []interface{}) {
		email, ok1 := stat[0].(string)
		uplink, ok2 := stat[1].(float64)
		downlink, ok3 := stat[2].(float64)
		// 错误数据直接跳过
		if ok1 && ok2 && ok3 == false {
			return
		}

		var u model.User
		if err := db.Where(&model.User{Email: email}).Select("id").First(&u).Error; err != nil {
			logger.Errorf("无法找到用户, 邮箱是: %v. 错误原因: %v", email, err.Error())
			return
		}
		if u.ID == 0 {
			return
		}
		s := model.Stat{
			User:     u.ID,
			Uplink:   uplink,
			Downlink: downlink,
		}
		if err := db.Create(&s).Error; err != nil {
			logger.Error("流量统计添加失败, 用户 ID: %v", s.User)
		}
	}
	for _, stat := range stats {
		go saveStat(stat)
	}
}

func statsHandler(c echo.Context) (err error) {
	var data struct {
		Stats [][]interface{} `json:"stats"`
	}
	if err = c.Bind(&data); err != nil {
		return
	}
	go saveStats(
		model.GetDB(c),
		c.Logger(),
		data.Stats,
	)
	return c.String(http.StatusOK, "received")
}
