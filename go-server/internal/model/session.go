package model

import "gorm.io/datatypes"

type Session struct {
	BaseModel
	Demographics datatypes.JSON `gorm:"type:jsonb" json:"demographics"`
	Fingerprint  string         `gorm:"type:varchar(255)" json:"fingerprint"`
	IsDuplicate  bool           `gorm:"type:boolean;default:true" json:"is_duplicate"`
}
