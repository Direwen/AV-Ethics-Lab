### SUBJECT DEMOGRAPHICS

- **Age Group:** {{.Demographic.AgeRange}}
- **Gender:** {{.Demographic.Gender}}
- **Driving Experience:** {{.Demographic.DrivingExperience}}

### DECISION HISTORY

The subject faced the following {{len .Responses}} scenarios:

{{range .Responses}}
---
**Scenario Context:**
{{.Narrative}}

**Environmental Factors:**

- Visibility: {{.Factors.Visibility}}
- Road Condition: {{.Factors.RoadCondition}}
- Speed: {{.Factors.Speed}}
- Brake Status: {{.Factors.BrakeStatus}}
- Tailgater Present: {{.Factors.HasTailgater}}
- Primary Entity at Risk: {{.Factors.PrimaryEntity}}
- Primary Entity Behavior: {{.Factors.PrimaryBehavior}}

**Subject's Response:**

- **Action Chosen:** {{index .RankedOptions 0}} (Top Choice)
- **Time Taken:** {{.ResponseTimeMs}}ms
- **Did they Timeout?** {{.IsTimeout}}
- **Did they Interact?** {{.HasInteracted}}
{{end}}

### INSTRUCTION

Analyze the data above. Ignore "Timeout" decisions when calculating their ethical preference, but note them as "Hesitation" in the summary if frequent. Generate the JSON profile now.