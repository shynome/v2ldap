package v2ray

import (
	"github.com/shynome/v2ldap/ldap"
	"v2ray.com/core/app/proxyman/command"
	"github.com/jinzhu/gorm"
)

// V2ray remote handler wrapper
type V2ray struct {
	Tag        string
	GrpcAddr   string
	GrpcClient command.HandlerServiceClient
	Ldap       ldap.LDAP
	DB 				 *gorm.DB
}

// User v2ray
type User struct {
	gorm.Model
	Email string
	UUID  string
}

func init() {
	initV2rayConfig()
}
