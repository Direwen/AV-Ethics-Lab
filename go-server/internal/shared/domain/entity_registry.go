package domain

type Entity struct {
	TypeID   string   `json:"type_id"`
	BaseName string   `json:"base_name"`
	Emoji    string   `json:"emoji"`
	Tags     []string `json:"tags"`
}

var EntityRegistry = map[string]Entity{

	// -- Ego --
	"vehicle_av": {
		TypeID: "vehicle_av", BaseName: "Autonomous Vehicle", Emoji: "🚕",
		Tags: []string{"ego", "vehicle", "agent"},
	},

	// --- STARS (Vulnerable / Moral Focus) ---
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

	// --- EXTRAS (Background Noise) ---
	"vehicle_car": {
		TypeID: "vehicle_car", BaseName: "Sedan", Emoji: "🚗",
		Tags: []string{"background", "vehicle"},
	},
	"vehicle_bus": {
		TypeID: "vehicle_bus", BaseName: "Bus", Emoji: "🚌",
		Tags: []string{"background", "vehicle", "large"},
	},
	"ped_adult": {
		TypeID: "ped_adult", BaseName: "Adult", Emoji: "🧍",
		Tags: []string{"background", "pedestrian"},
	},
	"obstacle_barrier": {
		TypeID: "obstacle_barrier", BaseName: "Barrier", Emoji: "🚧",
		Tags: []string{"background", "static"},
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
}

// POOLS (Cached lists for the Generator)
var (
	StarPool       = []string{"ped_child", "ped_elderly", "ped_doctor", "ped_criminal"}
	BackgroundPool = []string{"vehicle_car", "vehicle_bus", "ped_adult", "obstacle_barrier", "vehicle_sports_car", "vehicle_police", "vehicle_ambulance"}
)
