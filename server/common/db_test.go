package common

import (
	"testing"

	"github.com/shynome/v2ldap/v2ray"
)

func TestDB(t *testing.T) {
	var err error

	var db = GetDB()
	defer db.Close()
	db.AutoMigrate(&v2ray.User{})

	email := "string"

	if err = db.Create(&v2ray.User{Email: email, UUID: "string"}).Error; err != nil {
		t.Error(err)
		return
	}

	var u v2ray.User
	if err = db.First(&u, "email = ?", email).Error; err != nil {
		t.Error(err)
		return
	}

	if err = db.Unscoped().Delete(&u).Error; err != nil {
		t.Error(err)
		return
	}

	t.Log(u)
}
