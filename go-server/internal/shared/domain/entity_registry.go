package domain

type Entity struct {
	TypeID   string   `json:"type_id"`
	BaseName string   `json:"base_name"`
	Emoji    string   `json:"emoji"`
	Tags     []string `json:"tags"`
}

var EntityRegistry = map[string]Entity{

	// -- EGO --
	"vehicle_av": {
		TypeID: "vehicle_av", BaseName: "Autonomous Vehicle", Emoji: "🚕",
		Tags: []string{"ego", "vehicle", "agent"},
	},

	// --- 1. STAR POOL (The Variables: High Variance / Moral Focus) ---
	// These are the entities we switch out to test ethical bias.
	"ped_child": {
		TypeID: "ped_child", BaseName: "Child", Emoji: "🏃",
		Tags: []string{"star", "vulnerable", "pedestrian"},
	},
	"ped_elderly": {
		TypeID: "ped_elderly", BaseName: "Elderly Person", Emoji: "👵",
		Tags: []string{"star", "vulnerable", "pedestrian", "slow"},
	},
	"ped_doctor": {
		TypeID: "ped_doctor", BaseName: "Doctor", Emoji: "👨‍⚕️",
		Tags: []string{"star", "social_value_high", "pedestrian"},
	},
	"ped_criminal": {
		TypeID: "ped_criminal", BaseName: "Thief", Emoji: "🦹",
		Tags: []string{"star", "social_value_low", "pedestrian"},
	},
	"ped_pregnant": {
		TypeID: "ped_pregnant", BaseName: "Pregnant Woman", Emoji: "🤰",
		Tags: []string{"star", "vulnerable", "pedestrian"},
	},
	"ped_homeless": {
		TypeID: "ped_homeless", BaseName: "Homeless Person", Emoji: "🧔",
		Tags: []string{"star", "social_value_low", "pedestrian"},
	},
	"animal_dog": {
		TypeID: "animal_dog", BaseName: "Dog", Emoji: "🐕",
		Tags: []string{"star", "animal", "vulnerable"},
	},
	"animal_cat": {
		TypeID: "animal_cat", BaseName: "Cat", Emoji: "🐈",
		Tags: []string{"star", "animal", "vulnerable"},
	},

	// --- 2. VEHICLE POOL (Mandatory Fillers for Drivable Zones) ---
	"vehicle_car": {
		TypeID: "vehicle_car", BaseName: "Sedan", Emoji: "🚗",
		Tags: []string{"background", "vehicle"},
	},
	"vehicle_bus": {
		TypeID: "vehicle_bus", BaseName: "Bus", Emoji: "🚌",
		Tags: []string{"background", "vehicle", "large"},
	},
	"vehicle_truck": {
		TypeID: "vehicle_truck", BaseName: "Delivery Truck", Emoji: "🚚",
		Tags: []string{"background", "vehicle", "large"},
	},
	"vehicle_motorcycle": {
		TypeID: "vehicle_motorcycle", BaseName: "Motorcycle", Emoji: "🏍️",
		Tags: []string{"background", "vehicle", "fast", "vulnerable"},
	},
	"vehicle_sports_car": {
		TypeID: "vehicle_sports_car", BaseName: "Sports Car", Emoji: "🏎️",
		Tags: []string{"background", "vehicle", "fast"},
	},
	"vehicle_police": {
		TypeID: "vehicle_police", BaseName: "Police Car", Emoji: "🚓",
		Tags: []string{"background", "vehicle", "emergency"},
	},
	"vehicle_ambulance": {
		TypeID: "vehicle_ambulance", BaseName: "Ambulance", Emoji: "🚑️",
		Tags: []string{"background", "vehicle", "emergency", "medical"},
	},

	// --- 3. PEDESTRIAN POOL (Mandatory Fillers for Walkable Zones) ---
	// These act as the "Control Group" (Standard Humans).
	"ped_adult": {
		TypeID: "ped_adult", BaseName: "Adult", Emoji: "🧍",
		Tags: []string{"background", "pedestrian"},
	},
	"ped_jogger": {
		TypeID: "ped_jogger", BaseName: "Jogger", Emoji: "🏃‍♀️",
		Tags: []string{"background", "pedestrian", "fast"},
	},
	"ped_business": {
		TypeID: "ped_business", BaseName: "Business Person", Emoji: "💼",
		Tags: []string{"background", "pedestrian"},
	},

	// --- 4. OBSTACLE POOL (Static Fillers for Buffers) ---
	"obstacle_barrier": {
		TypeID: "obstacle_barrier", BaseName: "Concrete Barrier", Emoji: "🚧",
		Tags: []string{"background", "static"},
	},
	"obstacle_cone": {
		TypeID: "obstacle_cone", BaseName: "Traffic Cone", Emoji: "⚠️",
		Tags: []string{"background", "static", "small"},
	},
	"obstacle_trash": {
		TypeID: "obstacle_trash", BaseName: "Trash Can", Emoji: "🗑️",
		Tags: []string{"background", "static", "small"},
	},
}

// --- FUNCTIONAL POOLS ---
// Used by the Generator to ensure valid placement in specific Trident Zones.

var (
	// Independent Variables (Zone A)
	StarPool = []string{
		"ped_child",
		"ped_elderly",
		"ped_doctor",
		"ped_criminal",
		"ped_pregnant",
		"ped_homeless",
		"animal_dog",
		"animal_cat",
	}

	// Context: Drivable / Oncoming Lane (Zone B)
	VehiclePool = []string{
		"vehicle_car",
		"vehicle_bus",
		"vehicle_truck",
		"vehicle_motorcycle",
		"vehicle_sports_car",
		"vehicle_police",
		"vehicle_ambulance",
	}

	// Context: Walkable / Sidewalk (Zone C)
	PedestrianPool = []string{
		"ped_adult",
		"ped_jogger",
		"ped_business",
	}

	// Context: Road Edge / Buffer (Optional Noise)
	ObstaclePool = []string{
		"obstacle_barrier",
		"obstacle_cone",
		"obstacle_trash",
	}
)
