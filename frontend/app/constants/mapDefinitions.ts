import type { CellDefinition } from "~/types/simulation";

export const CELL_DEFINITIONS: Record<string, CellDefinition> = {
    // BUILDINGS
    "0": { name: "Roof", class: "bg-gray-600", isInteractive: true, surface: "building" },
    "1": { name: "Building Edge Top", class: "bg-gray-600 border-b-4 border-gray-900", isInteractive: true, surface: "building" },
    "2": { name: "Building Edge Bottom", class: "bg-gray-600 border-t-4 border-gray-900", isInteractive: true, surface: "building" },

    // SIDEWALKS
    "3": { name: "Sidewalk Top", class: "bg-gray-200 border-b border-gray-300", isInteractive: true, surface: "walkable" },
    "4": { name: "Sidewalk Bottom", class: "bg-gray-200 border-t border-gray-300", isInteractive: true, surface: "walkable" },
    "5": { name: "Sidewalk Corner Top-Right", class: "bg-gray-200 rounded-br-lg", isInteractive: true, surface: "walkable" },
    "6": { name: "Sidewalk Corner Bottom-Right", class: "bg-gray-200 rounded-tr-lg", isInteractive: true, surface: "walkable" },
    "7": { name: "Sidewalk Corner Top-Left", class: "bg-gray-200 rounded-bl-lg", isInteractive: true, surface: "walkable" },
    "8": { name: "Sidewalk Corner Bottom-Left", class: "bg-gray-200 rounded-tl-lg", isInteractive: true, surface: "walkable" },

    // ROADS
    "9": { name: "Asphalt Horizontal", class: "bg-gray-800", isInteractive: true, surface: "drivable" },
    "10": { name: "Asphalt Vertical", class: "bg-gray-800", isInteractive: true, surface: "drivable" },
    "11": { name: "Intersection Box", class: "bg-gray-800", isInteractive: true, surface: "drivable" },

    // ROAD MARKINGS
    "12": { name: "Yellow Line Dash", class: "yellow-dash-line-h", isInteractive: true, surface: "drivable" },
    "13": { name: "Double Yellow Horizontal", class: "bg-gray-800 border-y-2 border-yellow-400 bg-yellow-300", isInteractive: true, surface: "drivable" },
    "14": { name: "Double Yellow Vertical", class: "bg-gray-800 border-x-2 border-yellow-400 bg-yellow-300", isInteractive: true, surface: "drivable" },

    // CROSSWALKS
    "15": { name: "Crosswalk Vertical", class: "crosswalk-v", isInteractive: true, surface: "drivable" },
    "16": { name: "Crosswalk Horizontal", class: "crosswalk-h", isInteractive: true, surface: "drivable" },

    // SOLID YELLOW LINES
    "17": { name: "Yellow Line Vertical", class: "yellow-line-v", isInteractive: true, surface: "drivable" },
    "18": { name: "Yellow Line Horizontal", class: "yellow-line-h", isInteractive: true, surface: "drivable" },

    "19": { name: "Sidewalk Left", class: "bg-gray-200 border-r-4 border-gray-300", isInteractive: true, surface: "walkable" },
    "20": { name: "Sidewalk Right", class: "bg-gray-200 border-l-4 border-gray-300", isInteractive: true, surface: "walkable" },
}
