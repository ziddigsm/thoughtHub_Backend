package types

import (
	"time"


)

type Users struct {
    ID        int            `gorm:"primaryKey;autoIncrement" json:"-"`
    Mail      string         `gorm:"size:255;unique;not null" json:"mail"`
    Name      string         `gorm:"size:100;not null" json:"name"`
    Username  string         `gorm:"size:20" json:"username,omitempty"`
    IsActive  bool           `gorm:"default:true" json:"-"`
    CreatedOn time.Time `gorm:"autoCreateTime" json:"-"`
    UpdatedOn time.Time `gorm:"autoUpdateTime" json:"-"`
}

type Menu struct {
	ID int `gorm:"primaryKey;autoIncrement" json:"id"`
	Options string `gorm:"size:50" json:"option"`
	Is_active bool `gorm:"default:true" json:"is_active"`
	Is_navbar bool `gorm:"default:true" json:"is_navbar"`
	CreatedOn time.Time `gorm:"autoCreateTime" json:"-"`
    UpdatedOn time.Time `gorm:"autoUpdateTime" json:"-"`
}