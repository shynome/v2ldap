package common

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // init sqlite
)

// DB sqlite db
var DB *gorm.DB

// DBConnect db
func DBConnect() {
	var err error
	DB, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
}
