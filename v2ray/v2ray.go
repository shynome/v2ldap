package v2ray

import (
	"github.com/jinzhu/gorm"
	"v2ray.com/core"
	"v2ray.com/core/app/proxyman/command"
)

// V2ray remote handler wrapper must sync at first time
type V2ray struct {
	RemoteTag  string                       // remote v2ray tag
	RemoteGrpc string                       // grpc addr
	grpcClient command.HandlerServiceClient //
	config     *core.Config                 //
	DB         *gorm.DB                     // db for storage user uuid
}

// GetConfig expose v2ray config
func (v2 V2ray) GetConfig() *core.Config {
	return v2.config
}

// User v2ray
type User struct {
	gorm.Model
	Email string `gorm:"unique_index;not null;column:email"`
	UUID  string `gorm:"unique;not null;column:uuid"`
}

func init() {
}
