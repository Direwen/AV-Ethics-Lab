package dashboard

// PublicStats represents the aggregated public dashboard statistics
type PublicStats struct {
	CompletedSessions int64 `json:"completed_sessions"`

	CountriesRepresented int64 `json:"countries_represented"`

	LeastHarmfulOutcome *OutcomeDistribution `json:"least_harmful_outcome"`

	SelfPreservationEffect *TailgaterEffect `json:"self_preservation_effect"`

	EntityComplianceEffect *ComplianceEffect `json:"entity_compliance_effect"`

	DecisionTimeDistribution *TimeDistribution `json:"decision_time_distribution"`

	ArchetypeDistribution []ArchetypeCount `json:"archetype_distribution"`
}

// OutcomeDistribution represents the distribution of least harmful outcome rankings
type OutcomeDistribution struct {
	Maintain    int64   `json:"maintain"`
	SwerveLeft  int64   `json:"swerve_left"`
	SwerveRight int64   `json:"swerve_right"`
	Total       int64   `json:"total"`
	MaintainPct float64 `json:"maintain_pct"`
	LeftPct     float64 `json:"swerve_left_pct"`
	RightPct    float64 `json:"swerve_right_pct"`
}

// TailgaterEffect compares maintain-lane rankings with/without tailgater presence
type TailgaterEffect struct {
	WithTailgater    *EffectMetric `json:"with_tailgater"`
	WithoutTailgater *EffectMetric `json:"without_tailgater"`
}

// EffectMetric represents a percentage metric with sample size
type EffectMetric struct {
	MaintainCount int64   `json:"maintain_count"`
	TotalCount    int64   `json:"total_count"`
	Percentage    float64 `json:"percentage"`
}

// ComplianceEffect compares maintain rankings based on entity behavior (compliant vs violation)
type ComplianceEffect struct {
	Compliant *EffectMetric `json:"compliant"`
	Violation *EffectMetric `json:"violation"`
}

// TimeDistribution represents decision time buckets for cognitive load analysis
type TimeDistribution struct {
	Under2Sec    int64 `json:"under_2s"`
	Between2And4 int64 `json:"between_2s_4s"`
	Between4And6 int64 `json:"between_4s_6s"`
	Over6Sec     int64 `json:"over_6s"`
	Total        int64 `json:"total"`
}

// ArchetypeCount represents a single archetype with its occurrence count
type ArchetypeCount struct {
	Archetype string `json:"archetype"`
	Count     int64  `json:"count"`
}
