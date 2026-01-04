package domain

type Tile struct {
	Name          string         `json:"name"`
	IsInteractive bool           `json:"is_interactive"`
	Definition    TileDefinition `json:"definition"`
}

type TileDefinition struct {
	SurfaceType  SurfaceType `json:"surface_type"`
	DefaultUsage []string    `json:"default_usage"`
}

type SurfaceType string

const (
	SurfaceDrivable   SurfaceType = "drivable"   // Roads
	SurfaceWalkable   SurfaceType = "walkable"   // Sidewalks
	SurfaceBuilding   SurfaceType = "building"   // Walls/Roofs (No placement allowed)
	SurfaceRestricted SurfaceType = "restricted" // Road markings (No placement allowed)
)

const (
	UsageVehicle    = "vehicle"
	UsagePedestrian = "pedestrian"
)

var TileRegistry = map[int]Tile{
	// OBSTACLES (Buildings) (Cannot place anything here)
	0: {Name: "Roof", IsInteractive: false, Definition: TileDefinition{SurfaceType: SurfaceBuilding, DefaultUsage: []string{}}},
	1: {Name: "Building Edge Top", IsInteractive: false, Definition: TileDefinition{SurfaceType: SurfaceBuilding, DefaultUsage: []string{}}},
	2: {Name: "Building Edge Bottom", IsInteractive: false, Definition: TileDefinition{SurfaceType: SurfaceBuilding, DefaultUsage: []string{}}},

	// --- SIDEWALKS (Safe Zones) ---
	// Physics: Concrete
	// Rules: Pedestrians only
	3:  {Name: "Sidewalk Top", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceWalkable, DefaultUsage: []string{UsagePedestrian}}},
	4:  {Name: "Sidewalk Bottom", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceWalkable, DefaultUsage: []string{UsagePedestrian}}},
	5:  {Name: "Sidewalk Corner Top-Right", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceWalkable, DefaultUsage: []string{UsagePedestrian}}},
	6:  {Name: "Sidewalk Corner Bottom-Right", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceWalkable, DefaultUsage: []string{UsagePedestrian}}},
	7:  {Name: "Sidewalk Corner Top-Left", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceWalkable, DefaultUsage: []string{UsagePedestrian}}},
	8:  {Name: "Sidewalk Corner Bottom-Left", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceWalkable, DefaultUsage: []string{UsagePedestrian}}},
	19: {Name: "Sidewalk Left", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceWalkable, DefaultUsage: []string{UsagePedestrian}}},
	20: {Name: "Sidewalk Right", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceWalkable, DefaultUsage: []string{UsagePedestrian}}},

	// --- ROADS (Drivable) ---
	// Physics: Asphalt
	// Rules: Vehicles only (usually)
	9:  {Name: "Asphalt Horizontal", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceDrivable, DefaultUsage: []string{UsageVehicle}}},
	10: {Name: "Asphalt Vertical", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceDrivable, DefaultUsage: []string{UsageVehicle}}},
	11: {Name: "Intersection Box", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceDrivable, DefaultUsage: []string{UsageVehicle}}},

	// --- ROAD MARKINGS (Restricted - No placement allowed) ---
	// Physics: Asphalt (It's just paint on road)
	12: {Name: "Yellow Line Dash", IsInteractive: false, Definition: TileDefinition{SurfaceType: SurfaceRestricted, DefaultUsage: []string{}}},
	13: {Name: "Double Yellow Horizontal", IsInteractive: false, Definition: TileDefinition{SurfaceType: SurfaceRestricted, DefaultUsage: []string{}}},
	14: {Name: "Double Yellow Vertical", IsInteractive: false, Definition: TileDefinition{SurfaceType: SurfaceRestricted, DefaultUsage: []string{}}},
	17: {Name: "Yellow Line Vertical", IsInteractive: false, Definition: TileDefinition{SurfaceType: SurfaceRestricted, DefaultUsage: []string{}}},
	18: {Name: "Yellow Line Horizontal", IsInteractive: false, Definition: TileDefinition{SurfaceType: SurfaceRestricted, DefaultUsage: []string{}}},

	// --- CROSSWALKS (Shared Zones) ---
	// Physics: Asphalt (Cars can drive over it)
	// Rules: SHARED (Both Vehicles and Pedestrians belong here)
	15: {Name: "Crosswalk Vertical", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceDrivable, DefaultUsage: []string{UsageVehicle, UsagePedestrian}}},
	16: {Name: "Crosswalk Horizontal", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceDrivable, DefaultUsage: []string{UsageVehicle, UsagePedestrian}}},
}
