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
	SurfaceAsphalt  SurfaceType = "asphalt"  // Drivable road surface
	SurfaceConcrete SurfaceType = "concrete" // Walkable sidewalk surface
	SurfaceObstacle SurfaceType = "obstacle" // Walls/Roofs (No placement allowed)
)

const (
	UsageVehicle    = "vehicle"
	UsagePedestrian = "pedestrian"
)

var TileRegistry = map[int]Tile{
	// OBSTACLES (Buildings) (Cannot place anything here)
	0: {Name: "Roof", IsInteractive: false, Definition: TileDefinition{SurfaceType: SurfaceObstacle, DefaultUsage: []string{}}},
	1: {Name: "Building Edge Top", IsInteractive: false, Definition: TileDefinition{SurfaceType: SurfaceObstacle, DefaultUsage: []string{}}},
	2: {Name: "Building Edge Bottom", IsInteractive: false, Definition: TileDefinition{SurfaceType: SurfaceObstacle, DefaultUsage: []string{}}},

	// --- SIDEWALKS (Safe Zones) ---
	// Physics: Concrete
	// Rules: Pedestrians only
	3:  {Name: "Sidewalk Top", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceConcrete, DefaultUsage: []string{UsagePedestrian}}},
	4:  {Name: "Sidewalk Bottom", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceConcrete, DefaultUsage: []string{UsagePedestrian}}},
	5:  {Name: "Sidewalk Corner Top-Right", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceConcrete, DefaultUsage: []string{UsagePedestrian}}},
	6:  {Name: "Sidewalk Corner Bottom-Right", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceConcrete, DefaultUsage: []string{UsagePedestrian}}},
	7:  {Name: "Sidewalk Corner Top-Left", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceConcrete, DefaultUsage: []string{UsagePedestrian}}},
	8:  {Name: "Sidewalk Corner Bottom-Left", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceConcrete, DefaultUsage: []string{UsagePedestrian}}},
	19: {Name: "Sidewalk Left", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceConcrete, DefaultUsage: []string{UsagePedestrian}}},
	20: {Name: "Sidewalk Right", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceConcrete, DefaultUsage: []string{UsagePedestrian}}},

	// --- ROADS (Drivable) ---
	// Physics: Asphalt
	// Rules: Vehicles only (usually)
	9:  {Name: "Asphalt Horizontal", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceAsphalt, DefaultUsage: []string{UsageVehicle}}},
	10: {Name: "Asphalt Vertical", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceAsphalt, DefaultUsage: []string{UsageVehicle}}},
	11: {Name: "Intersection Box", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceAsphalt, DefaultUsage: []string{UsageVehicle}}},

	// --- ROAD MARKINGS (Still Roads) ---
	// Physics: Asphalt (It's just paint on road)
	12: {Name: "Yellow Line Dash", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceAsphalt, DefaultUsage: []string{UsageVehicle}}},
	13: {Name: "Double Yellow Horizontal", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceAsphalt, DefaultUsage: []string{UsageVehicle}}},
	14: {Name: "Double Yellow Vertical", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceAsphalt, DefaultUsage: []string{UsageVehicle}}},
	17: {Name: "Yellow Line Vertical", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceAsphalt, DefaultUsage: []string{UsageVehicle}}},
	18: {Name: "Yellow Line Horizontal", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceAsphalt, DefaultUsage: []string{UsageVehicle}}},

	// --- CROSSWALKS (Shared Zones) ---
	// Physics: Asphalt (Cars can drive over it)
	// Rules: SHARED (Both Vehicles and Pedestrians belong here)
	15: {Name: "Crosswalk Vertical", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceAsphalt, DefaultUsage: []string{UsageVehicle, UsagePedestrian}}},
	16: {Name: "Crosswalk Horizontal", IsInteractive: true, Definition: TileDefinition{SurfaceType: SurfaceAsphalt, DefaultUsage: []string{UsageVehicle, UsagePedestrian}}},
}

// TilesBySurface groups tile IDs by their surface type
var TilesBySurface map[SurfaceType][]int

func init() {
	TilesBySurface = make(map[SurfaceType][]int)
	for id, tile := range TileRegistry {
		TilesBySurface[tile.Definition.SurfaceType] = append(
			TilesBySurface[tile.Definition.SurfaceType], id,
		)
	}
}
