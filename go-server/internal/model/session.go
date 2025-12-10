package model

import "gorm.io/datatypes"

type Session struct {
	BaseModel
	Demographics datatypes.JSON `gorm:"type:jsonb" json:"demographics"`
	Fingerprint  string         `gorm:"type:varchar(255)" json:"fingerprint"`
}
