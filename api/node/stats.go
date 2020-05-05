package node

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/shynome/v2ldap/model"
)

func saveStats(db *gorm.DB, logger echo.Logger, stats [][]interface{}) {
	var saveStat = func(stat []interface{}) {
		email, uplink, downlink := stat[0].(string), stat[1].(float64), stat[2].(float64)
		var u model.User
		if err := db.Unscoped().Where("email = ?", email).Select("id").First(&u); err != nil {
			logger.Errorf("无法找到用户, 邮箱是: %v", email)
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
