package domain

import "context"

// Task represents different LLM task types
type LLMTask string

const (
	TaskScenario LLMTask = "scenario"
	TaskFeedback LLMTask = "feedback"
)

// A marker interface for all LLM clients
type Client interface {
	IsLLMClient()
}

// LLMClient defines the interface for LLM operations
type LLMClient interface {
	GenerateScenario(ctx context.Context, req ScenarioLLMRequest) (*ScenarioLLMResponse, error)
}

// Feedback LLM
type FeedbackLLMClient interface {
	GenerateFeedback(ctx context.Context, req FeedbackLLMRequest) (*FeedbackLLMResponse, error)
}

// LLMPool defines the interface for LLM pool with automatic retry
type LLMPool interface {
	Execute(task LLMTask, cb func(client Client) (any, error)) (any, error)
	Register(task LLMTask, prefix string)
}

// FeedbackLLMRequest is the request payload for feedback generation
type FeedbackLLMRequest struct {
	Demographic Demographic        `json:"demographic"`
	Responses   []EnrichedResponse `json:"responses"`
}

// FeedbackLLMResponse is the response from the LLM for feedback generation
type FeedbackLLMResponse struct {
	Archetype string `json:"archetype"`
	Summary   string `json:"summary"`
	KeyTrait  string `json:"key_trait"`
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
