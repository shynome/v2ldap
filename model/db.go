package model

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // init sqlite
	"github.com/labstack/echo/v4"
)

var db *gorm.DB

// Init db
func Init() (err error) {
	if db == nil {
		dbpath := os.Getenv("DB_PATH")
		if dbpath == "" {
			dbpath = "v2ldap.db"
		}
		if db, err = gorm.Open("sqlite3", dbpath); err != nil {
			panic(err)
		}
		if err = db.AutoMigrate(&User{}, &Stat{}).Error; err != nil {
			panic(err)
		}
	}
	return
}

// RegisterDB in echo middleware
func RegisterDB(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("db", db)
		return next(c)
	}
}

// GetDB from context
func GetDB(c echo.Context) *gorm.DB {
	return c.Get("db").(*gorm.DB)
}
