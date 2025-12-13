package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Response struct {
	BaseModel
	RankingOrder datatypes.JSON `gorm:"type:jsonb" json:"ranking_order"`
	ReactionTime int64          `gorm:"not null" json:"reaction_time"`
	ScenarioID   uuid.UUID      `gorm:"type:uuid,not null" json:"scenario_id"`
	Scenario     Scenario       `gorm:"foreignKey:ScenarioID;references:ID"`
}
