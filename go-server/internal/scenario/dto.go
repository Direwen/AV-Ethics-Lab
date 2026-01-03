package scenario

import (
	"github.com/direwen/go-server/internal/shared/domain"
)

type GenerateRequest struct {
	SessionID string `json:"session_id" binding:"required,uuid"`
}

type GetNextResponse struct {
	Narrative  string                 `json:"narrative"`
	Entities   []EnrichedEntity       `json:"entities"`
	Factors    domain.ScenarioFactors `json:"factors"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
	GridData   [][]int                `json:"grid_data"`
	LaneConfig domain.LaneConfigMap   `json:"lane_config"`
}

type EnrichedEntity struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"`
	Emoji    string             `json:"emoji"`
	Row      int                `json:"row"`
	Col      int                `json:"col"`
	Metadata EnrichedEntityMeta `json:"metadata"`
}

type EnrichedEntityMeta struct {
	IsStar      bool   `json:"is_star"`
	IsEgo       bool   `json:"is_ego"`
	IsViolation bool   `json:"is_violation"`
	Action      string `json:"action"`
	Orientation string `json:"orientation"`
}
