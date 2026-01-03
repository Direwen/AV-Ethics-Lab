package domain

import "context"

// LLMClient defines the interface for LLM operations
type LLMClient interface {
	GenerateScenario(ctx context.Context, req ScenarioLLMRequest) (*ScenarioLLMResponse, error)
}

// ScenarioLLMRequest is the request payload for scenario generation
type ScenarioLLMRequest struct {
	TemplateName    string          `json:"template_name"`
	GridDimensions  string          `json:"grid_dimensions"`
	Factors         ScenarioFactors `json:"factors"`
	WalkableCells   [][2]int        `json:"walkable_cells"`
	DrivableCells   [][2]int        `json:"drivable_cells"`
	BuildingCells   [][2]int        `json:"building_cells"`
	RestrictedCells [][2]int        `json:"restricted_cells"`
	LaneConfig      LaneConfigMap   `json:"lane_config"`
}

// ScenarioLLMResponse is the response from the LLM for scenario generation
type ScenarioLLMResponse struct {
	Verification string      `json:"_verification"`
	Narrative    string      `json:"narrative"`
	Entities     []RawEntity `json:"entities"`
}

// RawEntity represents an entity returned by the LLM
type RawEntity struct {
	Type     string        `json:"type"`
	Row      int           `json:"row"`
	Col      int           `json:"col"`
	Metadata RawEntityMeta `json:"metadata"`
}

// RawEntityMeta contains metadata for a raw entity
type RawEntityMeta struct {
	IsStar      bool   `json:"is_star"`
	IsEgo       bool   `json:"is_ego"`
	IsViolation bool   `json:"is_violation"`
	Action      string `json:"action"`
	Orientation string `json:"orientation"`
}
