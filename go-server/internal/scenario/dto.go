package scenario

import (
	"github.com/direwen/go-server/internal/shared/domain"
)

type GenerateRequest struct {
	SessionID string `json:"session_id" binding:"required,uuid"`
}

type GetNextResponse struct {
	Narrative      string                 `json:"narrative"`
	DilemmaOptions domain.DilemmaOptions  `json:"dilemma_options"`
	Entities       []EnrichedEntity       `json:"entities"`
	Factors        domain.ScenarioFactors `json:"factors"`
	Width          int                    `json:"width"`
	Height         int                    `json:"height"`
	GridData       [][]int                `json:"grid_data"`
	LaneConfig     domain.LaneConfigMap   `json:"lane_config"`
}

type EnrichedEntity struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Emoji    string            `json:"emoji"`
	Row      int               `json:"row"`
	Col      int               `json:"col"`
	Metadata domain.EntityMeta `json:"metadata"`
}
