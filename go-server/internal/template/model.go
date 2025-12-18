package template

import (
	"github.com/direwen/go-server/internal/shared/model"
	"gorm.io/datatypes"
)

type ContextTemplate struct {
	model.BaseModel
	Name     string         `gorm:"not null" json:"name"`
	Width    int            `gorm:"type:int;not null" json:"width"`
	Height   int            `gorm:"type:int;not null" json:"height"`
	GridData datatypes.JSON `gorm:"type:jsonb;not null" json:"grid_data"`
	Meta     datatypes.JSON `gorm:"type:jsonb" json:"meta"`
}
