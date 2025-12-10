package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Scenario struct {
	BaseModel
	ContextTemplateID  uuid.UUID       `gorm:"type:uuid;not null" json:"context_template_id"`
	ContextTemplate    ContextTemplate `gorm:"foreignKey:ContextTemplateID" json:"context_template"`
	EntityPlacements   datatypes.JSON  `gorm:"type:jsonb" json:"entity_placements"`
	EnvironmentFactors datatypes.JSON  `gorm:"type:jsonb" json:"environment_factors"`
}
