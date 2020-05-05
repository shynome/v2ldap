package model

import (
	"github.com/jinzhu/gorm"
)

// Stat user
type Stat struct {
	gorm.Model
	User     uint `gorm:"not null;column:user" json:"user"`
	Uplink   int64  `gorm:"not null;column:uplink" json:"uplink"`
	Downlink int64  `gorm:"not null;column:downlink;" json:"downlink"`
}
