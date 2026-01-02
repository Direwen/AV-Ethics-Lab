package domain

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
