package domain

import "context"

// LLMClient defines the interface for LLM operations
type LLMClient interface {
	GenerateScenario(ctx context.Context, req ScenarioLLMRequest) (*ScenarioLLMResponse, error)
}

// ScenarioLLMRequest is the request payload for scenario generation
type ScenarioLLMRequest struct {
	TemplateName   string          `json:"template_name"`
	GridDimensions string          `json:"grid_dimensions"`
	Factors        ScenarioFactors `json:"factors"`
	WalkableCells  [][2]int        `json:"walkable_cells"`
	DrivableCells  [][2]int        `json:"drivable_cells"`
	BuildingCells  [][2]int        `json:"building_cells"`
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
	IsViolation bool   `json:"is_violation"`
	Action      string `json:"action"`
}

// ScenarioFactors contains the factors used to generate a scenario
type ScenarioFactors struct {
	Visibility         string   `json:"visibility"`
	RoadCondition      string   `json:"road_condition"`
	Location           string   `json:"location"`
	BrakeStatus        string   `json:"brake_status"`
	Speed              string   `json:"speed"`
	HasTailgater       bool     `json:"has_tailgater"`
	PrimaryEntity      string   `json:"primary_entity"`
	PrimaryBehavior    string   `json:"primary_behavior"`
	BackgroundEntities []string `json:"background_entities"`
}
