package domain

// Visibility conditions
type Visibility string

const (
	VisibilityClear Visibility = "Clear"
	VisibilityFog   Visibility = "Fog"
	VisibilityNight Visibility = "Night"
	VisibilityRain  Visibility = "Rain"
)

// Road conditions
type RoadCondition string

const (
	RoadConditionDry RoadCondition = "Dry"
	RoadConditionWet RoadCondition = "Wet"
	RoadConditionIcy RoadCondition = "Icy"
)

// Location codes
type Location string

const (
	LocationUS Location = "US"
	LocationUK Location = "UK"
	LocationCN Location = "CN"
	LocationFR Location = "FR"
)

// Brake status
type BrakeStatus string

const (
	BrakeStatusActive BrakeStatus = "Active"
	BrakeStatusFailed BrakeStatus = "Failed"
	BrakeStatusFade   BrakeStatus = "Fade"
)

// Speed levels
type Speed string

const (
	SpeedLow    Speed = "Low"
	SpeedMedium Speed = "Medium"
	SpeedHigh   Speed = "High"
)

// Behavior types
type Behavior string

const (
	BehaviorViolation Behavior = "Violation"
	BehaviorCompliant Behavior = "Compliant"
)

var (
	Visibilities   = []Visibility{VisibilityClear, VisibilityFog, VisibilityNight, VisibilityRain}
	RoadConditions = []RoadCondition{RoadConditionDry, RoadConditionWet, RoadConditionIcy}
	Locations      = []Location{LocationUS, LocationUK, LocationCN, LocationFR}
	BrakeStatuses  = []BrakeStatus{BrakeStatusActive, BrakeStatusFailed, BrakeStatusFade}
	Speeds         = []Speed{SpeedLow, SpeedMedium, SpeedHigh}
)

// Direction constants for lane config
type Direction string

const (
	DirectionNorth Direction = "N"
	DirectionSouth Direction = "S"
	DirectionEast  Direction = "E"
	DirectionWest  Direction = "W"
)

// LaneConfigMap represents the parsed lane configuration
type LaneConfigMap map[Direction][][2]int
