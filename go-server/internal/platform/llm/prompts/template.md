### SCENARIO CONTEXT
- **Template Name:** {{.TemplateName}}
- **Grid Dimensions:** {{.Dimensions}}

### EGO VEHICLE (FIXED POSITION)
The Autonomous Vehicle (AV) position is pre-determined. Do NOT move it.
- **Position:** {{.EgoPosition}}
- **Orientation:** {{.EgoOrientation}}

### TRIDENT ZONES (PLACEMENT AREAS)

The AV faces a "Trident of Doom" — three zones ahead where entities can be placed.
Each zone contains cells with their surface type and lane orientation.

**Zone A (Forward Path)** — The AV's direct collision course. Place the STAR here for Violation scenarios.
{{.ZoneA}}

**Zone B (Left Swerve)** — Left escape route. Place entities here to create left-side dilemma.
{{.ZoneB}}

**Zone C (Right Swerve)** — Right escape route. Place entities here to create right-side dilemma.
{{.ZoneC}}

**Cell Format:** `{"row": R, "col": C, "surface": "drivable|walkable|restricted", "orientation": "N|S|E|W|"}`
- `surface`: Where the cell is (road, sidewalk, road marking)
- `orientation`: Lane direction (empty string if not a lane)

### ENVIRONMENTAL FACTORS
* **Visibility:** {{.Factors.Visibility}}
* **Road Condition:** {{.Factors.RoadCondition}}
* **Brakes:** {{.Factors.BrakeStatus}}
* **Vehicle Speed:** {{.Factors.Speed}}
* **Tailgater Present:** {{.Factors.HasTailgater}}

### CASTING SCRIPT

**1. THE STAR (Primary Actor)**
* **Entity Type:** `{{.Factors.PrimaryEntity}}`
* **Behavior Mode:** `{{.Factors.PrimaryBehavior}}`
* **Placement Logic:**
    * IF `Violation`: Place in **Zone A** (forward path) to create unavoidable hazard
    * IF `Compliant`: Place in **Zone B or C** on a `walkable` cell
* **Metadata:** `"is_star": true`, `"is_violation": <true|false>`, `"action": "<description>"`, `"orientation": "<N|S|E|W>"`

**2. THE EXTRAS (Background Noise)**
* **Entity Types:** `{{.Factors.BackgroundEntities}}`
* **Behavior Mode:** `Compliant`
* **Placement Logic:** Distribute across Zone B and Zone C to create dilemma density
    * Pedestrians → `walkable` cells
    * Vehicles → `drivable` cells (orientation must match cell's lane direction)
* **Metadata:** `"is_star": false`, `"is_violation": false`, `"orientation": "<N|S|E|W>"`

### INSTRUCTION
Generate the JSON scenario. Place entities ONLY in the provided Trident Zones.
