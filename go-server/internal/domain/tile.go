package domain

import (
	"github.com/direwen/go-server/internal/dto"
)

var TileRegistry = map[int]dto.Tile{
	// BUILDINGS
	0: {Name: "Roof", IsInteractive: false},
	1: {Name: "Building Edge Top", IsInteractive: false},
	2: {Name: "Building Edge Bottom", IsInteractive: false},

	// SIDEWALKS
	3: {Name: "Sidewalk Top", IsInteractive: true},
	4: {Name: "Sidewalk Bottom", IsInteractive: true},
	5: {Name: "Sidewalk Corner Top-Right", IsInteractive: true},
	6: {Name: "Sidewalk Corner Bottom-Right", IsInteractive: true},
	7: {Name: "Sidewalk Corner Top-Left", IsInteractive: true},
	8: {Name: "Sidewalk Corner Bottom-Left", IsInteractive: true},

	// ROADS
	9:  {Name: "Asphalt Horizontal", IsInteractive: true},
	10: {Name: "Asphalt Vertical", IsInteractive: true},
	11: {Name: "Intersection Box", IsInteractive: true},

	// ROAD MARKINGS
	12: {Name: "Yellow Line Dash", IsInteractive: true},
	13: {Name: "Double Yellow Horizontal", IsInteractive: true},
	14: {Name: "Double Yellow Vertical", IsInteractive: true},

	// CROSSWALKS
	15: {Name: "Crosswalk Vertical", IsInteractive: true},
	16: {Name: "Crosswalk Horizontal", IsInteractive: true},

	// SOLID YELLOW LINES
	17: {Name: "Yellow Line Vertical", IsInteractive: true},
	18: {Name: "Yellow Line Horizontal", IsInteractive: true},

	// MORE SIDEWALKS
	19: {Name: "Sidewalk Left", IsInteractive: true},
	20: {Name: "Sidewalk Right", IsInteractive: true},
}
