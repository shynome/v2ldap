package common

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // init sqlite
	"github.com/shynome/v2ldap/v2ray"
)

// DB sqlite db
var DB *gorm.DB

// GetDB return DB
func GetDB() *gorm.DB {
	var err error
	if DB == nil {
		if DB, err = gorm.Open("sqlite3", "test.db"); err != nil {
			panic(err)
		}
		if err = DB.AutoMigrate(&v2ray.User{}).Error; err != nil {
			panic(err)
		}
	}
	return DB
}
