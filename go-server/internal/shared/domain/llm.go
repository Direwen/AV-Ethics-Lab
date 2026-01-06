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

	// Ego AV (pre-determined position)
	EgoPosition    Coordinate `json:"ego_position"`
	EgoOrientation Direction  `json:"ego_orientation"`

	// Trident zones (enriched with surface/orientation)
	TridentZones TridentZones `json:"trident_zones"`
}

// ScenarioLLMResponse is the response from the LLM for scenario generation
type ScenarioLLMResponse struct {
	Verification   string         `json:"_verification"`
	Narrative      string         `json:"narrative"`
	DilemmaOptions DilemmaOptions `json:"dilemma_options"`
	Entities       []RawEntity    `json:"entities"`
}

// DilemmaOptions contains the text for the 3 user action buttons
type DilemmaOptions struct {
	Maintain    string `json:"maintain"`
	SwerveLeft  string `json:"swerve_left"`
	SwerveRight string `json:"swerve_right"`
}

// EntityMeta contains metadata for an entity (shared between LLM and enriched entities)
type EntityMeta struct {
	IsStar      bool   `json:"is_star"`
	IsEgo       bool   `json:"is_ego"`
	IsViolation bool   `json:"is_violation"`
	Action      string `json:"action"`
	Orientation string `json:"orientation"`
}

// RawEntity represents an entity returned by the LLM
type RawEntity struct {
	Type     string     `json:"type"`
	Row      int        `json:"row"`
	Col      int        `json:"col"`
	Metadata EntityMeta `json:"metadata"`
}
