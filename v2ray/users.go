package v2ray

// Users in memory
var Users []string

func initMemoryUsers() {
	var err error
	Users, err = Ldap.GetUsers()
	if err != nil {
		panic(err)
	}
}
