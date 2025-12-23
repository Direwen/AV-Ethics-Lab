### SCENARIO CONTEXT
* **Template Name:** {{.TemplateName}}
* **Grid Dimensions:** {{.Dimensions}}
* **Map Layout:** {{.GridData}}

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
* **Metadata Requirement:** Set `"is_star": true`, `"is_violation": {{if eq .Factors.PrimaryBehavior "Violation"}}true{{else}}false{{end}}`.

**2. THE EXTRAS (Background Noise)**
* **Entity Types:** `{{.Factors.BackgroundEntities}}`
* **Behavior Mode:** `Compliant`
* **Placement Logic:** Distribute these to create realistic scene density. They must strictly follow traffic rules (Sidewalks for peds, Lanes for cars).
* **Metadata Requirement:** Set `"is_star": false`, `"is_violation": false`.

### INSTRUCTION
Generate the JSON scenario now based on the **Map Layout** and **Casting Script** above.