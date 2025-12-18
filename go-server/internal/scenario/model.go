package scenario

import (
	"github.com/direwen/go-server/internal/session"
	"github.com/direwen/go-server/internal/shared/model"
	"github.com/direwen/go-server/internal/template"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Scenario struct {
	model.BaseModel
	ContextTemplateID  uuid.UUID                `gorm:"type:uuid;not null" json:"context_template_id"`
	ContextTemplate    template.ContextTemplate `gorm:"foreignKey:ContextTemplateID" json:"context_template"`
	SessionID          uuid.UUID                `gorm:"type:uuid;not null" json:"session_id"`
	Session            session.Session          `gorm:"foreignKey:SessionID" json:"session"`
	EntityPlacements   datatypes.JSON           `gorm:"type:jsonb" json:"entity_placements"`
	EnvironmentFactors datatypes.JSON           `gorm:"type:jsonb" json:"environment_factors"`
}
