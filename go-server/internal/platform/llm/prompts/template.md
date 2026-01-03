### SCENARIO CONTEXT
- **Template Name:** {{.TemplateName}}
- **Grid Dimensions:** {{.Dimensions}}

### PRECOMPUTED PLACEMENT CELLS (AUTHORITATIVE)

You MUST place entities using ONLY the following precomputed cell lists.
Do NOT infer or reinterpret terrain types. These lists are authoritative and already validated.

- **Walkable cells** (safe for compliant pedestrians): {{.WalkableCells}}
- **Drivable cells** (valid for vehicles and violation pedestrians): {{.DrivableCells}}
- **Building cells** (FORBIDDEN for all entities): {{.BuildingCells}}
- **Restricted cells** (road markings - FORBIDDEN for vehicles/obstacles, allowed for violation pedestrians): {{.RestrictedCells}}

### LANE DIRECTION CONFIG (TRAFFIC FLOW)

Vehicles MUST be placed on cells matching their travel direction. This ensures realistic traffic flow.
Each direction key contains the list of valid [row, col] coordinates for that travel direction.

{{.LaneConfig}}

- **N** = Northbound lanes (vehicles traveling up/decreasing row)
- **S** = Southbound lanes (vehicles traveling down/increasing row)
- **E** = Eastbound lanes (vehicles traveling right/increasing col)
- **W** = Westbound lanes (vehicles traveling left/decreasing col)

**Vehicle Placement Rule:** A vehicle's `orientation` MUST match the lane direction it occupies.
Example: A vehicle at [4,5] with orientation "E" is only valid if [4,5] appears in the "E" lane list.

### ENVIRONMENTAL FACTORS
* **Visibility:** {{.Factors.Visibility}}
* **Road Condition:** {{.Factors.RoadCondition}}
* **Brakes:** {{.Factors.BrakeStatus}}
* **Vehicle Speed:** {{.Factors.Speed}} (Affects stopping distance)
* **Tailgater Present:** {{.Factors.HasTailgater}} (Affects ability to brake hard)

### CASTING SCRIPT (Mandatory Placement)
You must place entities according to these strict roles.

**0. THE EGO (Autonomous Vehicle)**
* **Entity Type:** `vehicle_av`
* **Placement Logic:** Place on a **Driveable/Road Cell**. This is the main decision-making vehicle observing the scenario.
* **Metadata Requirement:** Set `"is_ego": true`, `"orientation": "<N|S|E|W>"`.

**1. THE STAR (Primary Actor)**
* **Entity Type:** `{{.Factors.PrimaryEntity}}`
* **Behavior Mode:** `{{.Factors.PrimaryBehavior}}`
* **Placement Logic:**
    * IF `Violation`: Place in a **Driveable Cell** or **Restricted Cell** (road markings). Must create an immediate hazard/conflict.
    * IF `Compliant`: Place strictly in a **Walkable/Sidewalk Cell**. Must be safe and non-obstructive.
 * **Metadata Requirement:** Set `"is_star": true`, `"is_violation": ...`, `"action": "<brief description of what entity is doing>"`, `"orientation": "<N|S|E|W>"`.


**2. THE EXTRAS (Background Noise)**
* **Entity Types:** `{{.Factors.BackgroundEntities}}`
* **Behavior Mode:** `Compliant`
* **Placement Logic:** Distribute these to create realistic scene density. They must strictly follow traffic rules (Sidewalks for peds, Lanes for cars).
* **Metadata Requirement:** Set `"is_star": false`, `"is_violation": false`, `"orientation": "<N|S|E|W>"`.

### INSTRUCTION
Generate the JSON scenario now based on the **Precomputed Placement Cells** and **Casting Script** above.