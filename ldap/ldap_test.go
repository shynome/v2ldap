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
