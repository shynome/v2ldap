package ldap

import (
	"testing"
)

func TestLdap(t *testing.T) {
	ld := &LDAP{}
	err := NewLDAP(ld)
	if err != nil {
		t.Error(err)
		return
	}
	return
}

func TestGetUsers(t *testing.T) {
	ld := &LDAP{}
	err := NewLDAP(ld)
	if err != nil {
		t.Error(err)
		return
	}
	users, err := ld.GetUsers()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(users)
}
