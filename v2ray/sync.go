package v2ray

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/thoas/go-funk"
)

var syncMux = sync.Mutex{}

// SyncResponse of sync func
type SyncResponse struct {
	Add    []string
	Delete []string
}

func emailToUser(email string) User {
	return User{Email: email}
}

// Sync ldap user to here
// TODO: add test v2ray grpc manager
func (v2 V2ray) Sync(ldapUsers []string, confirm bool) (resp SyncResponse, err error) {

	usreAdd := []string{}
	userExist := []User{}
	userDelete := []string{}

	if err = v2.DB.Where("email in (?)", ldapUsers).Select([]string{"email"}).Find(&userExist).Error; err != nil {
		return
	}
	var userExistMap = map[string]bool{}
	for _, user := range userExist {
		userExistMap[user.Email] = true
	}
	usreAdd, ok := (funk.Filter(ldapUsers, func(email string) bool {
		return userExistMap[email] == false
	})).([]string)
	if ok == false {
		err = fmt.Errorf("ldap user filter fail")
		return
	}

	var userDeleteFromDB []User
	if err = v2.DB.Not("email", ldapUsers).Find(&userDeleteFromDB).Error; err != nil {
		return
	}

	userDelete = (funk.Map(userDeleteFromDB, func(user User) string {
		return user.Email
	})).([]string)
	resp = SyncResponse{
		Delete: userDelete,
		Add:    usreAdd,
	}

	if confirm == false {
		return
	}
	syncMux.Lock()
	defer syncMux.Unlock()

	userToAdd := funk.Map(usreAdd, func(email string) User {
		var id string
		// 如果用户已存在的话使用已存在的用户 UUID
		var u User
		if v2.DB.Unscoped().First(&u, "email = ?", email).Error == nil {
			id = u.UUID
		} else {
			id = uuid.New().String()
		}
		return User{
			Email: email,
			UUID:  id,
		}
	}).([]User)
	userToDelete := funk.Map(userDelete, emailToUser).([]User)

	// update to db
	if err = v2.DB.Where("email in (?)", userDelete).Delete(&User{}).Error; err != nil {
		return
	}
	errs := v2.loopUsers(func(user User) (err error) {
		var count int64
		if err = v2.DB.Model(&User{}).Unscoped().Where("email = ?", user.Email).Count(&count).Error; err != nil {
			return
		}
		// active exist user
		if count != 0 {
			if err = v2.DB.Model(&User{}).Unscoped().Where("email = ?", user.Email).Update("deleted_at", nil).Error; err != nil {
				return
			}
			return
		}
		if err = v2.DB.Create(&user).Error; err != nil {
			err = fmt.Errorf("add user %v throw err: %v", user.Email, err.Error())
			return
		}
		return
	})(userToAdd)
	if len(errs) != 0 {
		err = fmt.Errorf("add users has errors : %v", errs)
		return
	}

	// update to remtoe v2ray
	if v2.RemoteGrpc != "" && v2.RemoteTag != "" {
		v2.RemoveUsers(userToDelete)
		v2.AddUsers(userToAdd)
	}

	// reset config to nil
	v2.config = nil

	return
}
