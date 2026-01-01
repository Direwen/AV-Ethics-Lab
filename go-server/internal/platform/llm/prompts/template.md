### SCENARIO CONTEXT
- **Template Name:** {{.TemplateName}}
- **Grid Dimensions:** {{.Dimensions}}

### PRECOMPUTED PLACEMENT CELLS (AUTHORITATIVE)

You MUST place entities using ONLY the following precomputed cell lists.
Do NOT infer or reinterpret terrain types. These lists are authoritative and already validated.

- **Walkable cells** (safe for compliant pedestrians): {{.WalkableCells}}
- **Drivable cells** (valid for vehicles and violation pedestrians): {{.DrivableCells}}
- **Building cells** (FORBIDDEN for all entities): {{.BuildingCells}}

### ENVIRONMENTAL FACTORS
* **Visibility:** {{.Factors.Visibility}}
* **Road Condition:** {{.Factors.RoadCondition}}
* **Brakes:** {{.Factors.BrakeStatus}}
* **Vehicle Speed:** {{.Factors.Speed}} (Affects stopping distance)
* **Tailgater Present:** {{.Factors.HasTailgater}} (Affects ability to brake hard)

### CASTING SCRIPT (Mandatory Placement)
You must place entities according to these strict roles.

**1. THE STAR (Primary Actor)**
* **Entity Type:** `{{.Factors.PrimaryEntity}}`
* **Behavior Mode:** `{{.Factors.PrimaryBehavior}}`
* **Placement Logic:**
    * IF `Violation`: Place strictly in a **Driveable/Road Cell**. Must create an immediate hazard/conflict.
    * IF `Compliant`: Place strictly in a **Walkable/Sidewalk Cell**. Must be safe and non-obstructive.
 * **Metadata Requirement:** Set `"is_star": true`, `"is_violation": ...`, `"action": "<brief description of what entity is doing>"`.


**2. THE EXTRAS (Background Noise)**
* **Entity Types:** `{{.Factors.BackgroundEntities}}`
* **Behavior Mode:** `Compliant`
* **Placement Logic:** Distribute these to create realistic scene density. They must strictly follow traffic rules (Sidewalks for peds, Lanes for cars).
* **Metadata Requirement:** Set `"is_star": false`, `"is_violation": false`.

### INSTRUCTION
Generate the JSON scenario now based on the **Precomputed Placement Cells** and **Casting Script** above.