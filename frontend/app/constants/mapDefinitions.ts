export const CELL_DEFINITIONS: Record<string, { name: string; class: string }> = {
    // BUILDINGS
    "0": { name: "Roof", class: "bg-gray-600" },
    "1": { name: "Building Edge Top", class: "bg-gray-600 border-b-4 border-gray-900" },
    "2": { name: "Building Edge Bottom", class: "bg-gray-600 border-t-4 border-gray-900" },

    // SIDEWALKS
    "3": { name: "Sidewalk Top", class: "bg-gray-200 border-b border-gray-300" },
    "4": { name: "Sidewalk Bottom", class: "bg-gray-200 border-t border-gray-300" },
    "5": { name: "Sidewalk Corner Top-Right", class: "bg-gray-200 rounded-br-lg" },
    "6": { name: "Sidewalk Corner Bottom-Right", class: "bg-gray-200 rounded-tr-lg" },
    "7": { name: "Sidewalk Corner Top-Left", class: "bg-gray-200 rounded-bl-lg" },
    "8": { name: "Sidewalk Corner Bottom-Left", class: "bg-gray-200 rounded-tl-lg" },

    // ROADS
    "9": { name: "Asphalt Horizontal", class: "bg-gray-800" },
    "10": { name: "Asphalt Vertical", class: "bg-gray-800" },
    "11": { name: "Intersection Box", class: "bg-gray-800" },

    // ROAD MARKINGS
    "12": { name: "Yellow Line Dash", class: "yellow-dash-line-h" },
    "13": { name: "Double Yellow Horizontal", class: "bg-gray-800 border-y-2 border-yellow-400" },
    "14": { name: "Double Yellow Vertical", class: "bg-gray-800 border-x-2 border-yellow-400" },

    // CROSSWALKS
    "15": { name: "Crosswalk Vertical", class: "crosswalk-v" },
    "16": { name: "Crosswalk Horizontal", class: "crosswalk-h" },
}
