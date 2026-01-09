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

type Coordinate struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// EnrichedCoordinate adds surface and orientation info to a coordinate
type EnrichedCoordinate struct {
	Coordinate
	Surface     SurfaceType `json:"surface"`
	Orientation Direction   `json:"orientation"` // lane direction, empty if not a lane
}

type TridentSpawn struct {
	Coordinate
	Orientation Direction `json:"orientation"`
}

type TridentZone struct {
	Coordinates []EnrichedCoordinate `json:"coordinates"`
}

type TridentZones struct {
	ZoneA TridentZone `json:"zone_a"`
	ZoneB TridentZone `json:"zone_b"`
	ZoneC TridentZone `json:"zone_c"`
}

type Demographic struct {
	AgeRange          int    `json:"age_range"`
	Gender            int    `json:"gender"`
	Country           string `json:"country"`
	Occupation        string `json:"occupation"`
	DrivingExperience int    `json:"driving_experience"`
}

type EnrichedResponse struct {
	Narrative      string          `json:"narrative"`
	Factors        ScenarioFactors `json:"factors"`
	RankedOptions  []string        `json:"ranked_options"`
	HasInteracted  bool            `json:"has_interacted"`
	ResponseTimeMs int64           `json:"response_time_ms"`
	IsTimeout      bool            `json:"is_timeout"`
}
