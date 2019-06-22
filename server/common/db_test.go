package common

import (
	"testing"

	"github.com/shynome/v2ldap/v2ray"
)

func TestDB(t *testing.T) {
	var err error

	DBConnect()
	defer DB.Close()
	DB.AutoMigrate(&v2ray.User{})

	email := "string"

	if err = DB.Create(&v2ray.User{Email: email, UUID: "string"}).Error; err != nil {
		t.Error(err)
		return
	}

	var u v2ray.User
	if err = DB.First(&u, "email = ?", email).Error; err != nil {
		t.Error(err)
		return
	}

	if err = DB.Unscoped().Delete(&u, "email = ?", email).Error; err != nil {
		t.Error(err)
		return
	}

	t.Log(u)
}
