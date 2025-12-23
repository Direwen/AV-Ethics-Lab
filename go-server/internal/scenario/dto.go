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
