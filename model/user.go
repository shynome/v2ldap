package model

import (
	"github.com/jinzhu/gorm"
)

// User v2ray
type User struct {
	gorm.Model
	// 表示这个用户是第几版程序生成的, 第0版的话用户是直接放到配置里的, 第1版是走接口验证的了
	Version  int    `gorm:"column:version;" json:"version"`
	Email    string `gorm:"unique_index;not null;column:email" json:"email"`
	UUID     string `gorm:"unique;not null;column:uuid" json:"uuid"`
	Disabled bool   `gorm:"column:disabled" json:"disabled"`
	Remark   string `gorm:"column:remark;" json:"remark"`
}
