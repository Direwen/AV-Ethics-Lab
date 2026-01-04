### SCENARIO CONTEXT
- **Template:** {{.TemplateName}}
- **Grid:** {{.Dimensions}}
- **Factors:** {{.Factors.Visibility}}, {{.Factors.RoadCondition}}, Brakes: {{.Factors.BrakeStatus}}

### EGO VEHICLE
- **Position:** {{.EgoPosition}} (Fixed)

### TRIDENT ZONES (Available Slots)
These are the ONLY valid coordinates for placement.

**Zone A (Forward - The Threat):**
{{.ZoneA}}

**Zone B (Left - The Swerve):**
{{.ZoneB}}

**Zone C (Right - The Swerve):**
{{.ZoneC}}

### CASTING CALL

**1. THE STAR (The Independent Variable)**
* **Entity:** `{{.Factors.PrimaryEntity}}`
* **Behavior:** `{{.Factors.PrimaryBehavior}}`
* **Instructions:**
    * IF `Violation`: Place in **Zone A** as the immediate threat.
    * IF `Compliant`: Place in **Zone B or C** (whichever matches the entity type, e.g., Peds on Sidewalks).

**2. THE EXTRAS (The Filling Kit)**
* **Available:** `{{.Factors.BackgroundEntities}}`
* **Instructions:**
    * **STEP 1:** Check Zone A. Is the Star there?
        * YES: Distribute extras to B and C.
        * NO: **You MUST pick a Vehicle or Pedestrian from this list and place it in Zone A** to create a collision course.
    * **STEP 2:** Fill remaining zones (B and C) so the user has no "safe" option.

### GENERATION TASK
Generate the JSON output. Ensure **Zone A is populated**. Generate 3 distinct "Action Strings" for the UI based on the entities you placed.
