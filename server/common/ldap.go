package common

import (
	"github.com/shynome/v2ldap/ldap"
)

// Ldap instance
var Ldap = &ldap.LDAP{}

func initLdap() {
	if Ldap.BindDN != "" {
		return
	}
	ldap.NewLDAP(Ldap)
}
