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
	APIPort    uint32                       // remote v2ray api port
	WsPort     uint32                       // ws port
	WsPath     string                       // ws path
	config     *core.Config                 //
	DB         *gorm.DB                     // db for storage user uuid
	VNEXT      string                       // 如果有 VNEXT 将会替换默认 outbound
	SocksPort  uint32                       // 如果有值则暴露一个无需认证的 socks 端口
}

// User v2ray
type User struct {
	gorm.Model
	Email string `gorm:"unique_index;not null;column:email"`
	UUID  string `gorm:"unique;not null;column:uuid"`
}

func init() {
}
