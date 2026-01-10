package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Response struct {
	BaseModel
	ScenarioID     uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"scenario_id"`
	Scenario       *Scenario      `gorm:"foreignKey:ScenarioID" json:"-"`
	RankingOrder   datatypes.JSON `gorm:"type:jsonb" json:"ranking_order"`
	IsTimeout      bool           `gorm:"type:bool;default:false;not null" json:"is_timeout"`
	HasInteracted  bool           `gorm:"type:bool;default:true;not null" json:"has_interacted"`
	ResponseTimeMs int64          `gorm:"not null" json:"response_time_ms"`
}
