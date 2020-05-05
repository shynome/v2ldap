package model

import (
	"github.com/jinzhu/gorm"
)

// Stat user
type Stat struct {
	gorm.Model
	User     string  `gorm:"not null;column:user" json:"user"`
	Uplink   float64 `gorm:"not null;column:uplink" json:"uplink"`
	Downlink float64 `gorm:"not null;column:downlink;" json:"downlink"`
}
