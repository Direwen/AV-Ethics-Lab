package scenario

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type GenerateRequest struct {
	SessionID string `json:"session_id" binding:"required,uuid"`
}

type ScenarioResponse struct {
	ID        uuid.UUID      `json:"id"`
	Narrative string         `json:"narrative"`
	Entities  datatypes.JSON `json:"entities"`
	Factors   datatypes.JSON `json:"factors"`
	Width     int            `json:"width"`
	Height    int            `json:"height"`
	GridData  [][]int        `json:"grid_data"`
	Meta      datatypes.JSON `json:"meta"`
}

type GetNextResponse struct {
	Narrative string         `json:"narrative"`
	Entities  datatypes.JSON `json:"entities"`
	Factors   datatypes.JSON `json:"factors"`
	Width     int            `json:"width"`
	Height    int            `json:"height"`
	GridData  datatypes.JSON `json:"grid_data"`
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
