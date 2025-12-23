package response

import (
	"github.com/direwen/go-server/internal/shared/model"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Response struct {
	model.BaseModel
	RankingOrder datatypes.JSON `gorm:"type:jsonb" json:"ranking_order"`
	ReactionTime int64          `gorm:"not null" json:"reaction_time"`
	ScenarioID   uuid.UUID      `gorm:"type:uuid;not null" json:"scenario_id"`
}
