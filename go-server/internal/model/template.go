package model

import (
	"gorm.io/datatypes"
)

type ContextTemplate struct {
	BaseModel
	Name     string         `gorm:"not null" json:"name"`
	BaseGrid datatypes.JSON `gorm:"type:jsonb" json:"base_grid"`
}
