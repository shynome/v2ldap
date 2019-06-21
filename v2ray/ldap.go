package v2ray

import (
	"github.com/shynome/v2ldap/ldap"
)

// Ldap instance for vpn
var Ldap = &ldap.LDAP{}

func initLdapInstance() {
	if err := ldap.NewLDAP(Ldap); err != nil {
		panic(err)
	}
}
