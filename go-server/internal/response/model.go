package response

import (
	"github.com/direwen/go-server/internal/shared/model"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Response struct {
	model.BaseModel
	RankingOrder   datatypes.JSON `gorm:"type:jsonb" json:"ranking_order"`
	ScenarioID     uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"scenario_id"`
	IsTimeout      bool           `gorm:"type:bool;default:false;not null;" json:"is_timeout"`
	HasInteracted  bool           `gorm:"type:bool;default:true;not null;" json:"has_interacted"`
	ResponseTimeMs int64          `gorm:"not null" json:"response_time_ms"`
}
