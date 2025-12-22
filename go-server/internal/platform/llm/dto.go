package llm

import "github.com/direwen/go-server/internal/scenario"

type ScenarioRequest struct {
	TemplateName   string                   `json:"template_name"`
	GridDimensions string                   `json:"grid_dimensions"`
	GridData       [][]int                  `json:"grid_data"`
	Factors        scenario.ScenarioFactors `json:"factors"`
}

type ScenarioResponse struct {
	Verification string      `json:"_verification"`
	Narrative    string      `json:"narrative"`
	Entities     []RawEntity `json:"entities"`
}

type RawEntity struct {
	Type     string        `json:"type"`
	Row      int           `json:"row"`
	Col      int           `json:"col"`
	Metadata RawEntityMeta `json:"metadata"`
}

type RawEntityMeta struct {
	IsStar      bool   `json:"is_star"`
	IsViolation bool   `json:"is_violation"`
	Action      string `json:"action"`
}
