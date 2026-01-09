package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Scenario struct {
	BaseModel
	ContextTemplateID uuid.UUID        `gorm:"type:uuid;not null" json:"context_template_id"`
	ContextTemplate   *ContextTemplate `gorm:"foreignKey:ContextTemplateID" json:"context_template,omitempty"`
	SessionID         uuid.UUID        `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"session_id"`
	Session           *Session         `gorm:"foreignKey:SessionID" json:"-"`
	Entities          datatypes.JSON   `gorm:"type:jsonb" json:"entities"`
	Factors           datatypes.JSON   `gorm:"type:jsonb" json:"factors"`
	DilemmaOptions    datatypes.JSON   `gorm:"type:jsonb" json:"dilemma_options"`
	Narrative         string           `gorm:"type:text" json:"narrative"`
	TridentSpawn      datatypes.JSON   `gorm:"type:jsonb" json:"trident_spawn"`
	StartedAt         *time.Time       `gorm:"type:timestamp" json:"started_at"`
	// Relationship
	Response *Response `gorm:"foreignKey:ScenarioID" json:"response,omitempty"`
}
