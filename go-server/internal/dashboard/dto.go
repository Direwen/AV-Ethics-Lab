package dashboard

type PublicStats struct {
	CompletedSessions        int64                   `json:"completed_sessions"`
	CountriesRepresented     int64                   `json:"countries_represented"`
	LeastHarmfulOutcome      *OutcomeDistribution    `json:"least_harmful_outcome"`
	SelfPreservationEffect   *TailgaterEffect        `json:"self_preservation_effect"`
	EntityComplianceEffect   *ComplianceEffect       `json:"entity_compliance_effect"`
	DecisionTimeDistribution []TimeDistributionPoint `json:"decision_time_distribution"`
	ArchetypeDistribution    []ArchetypeCount        `json:"archetype_distribution"`
}

type OutcomeDistribution struct {
	Maintain    int64   `json:"maintain" gorm:"column:maintain"`
	SwerveLeft  int64   `json:"swerve_left" gorm:"column:swerve_left"`
	SwerveRight int64   `json:"swerve_right" gorm:"column:swerve_right"`
	Total       int64   `json:"total" gorm:"column:total"`
	MaintainPct float64 `json:"maintain_pct" gorm:"column:maintain_pct"`
	LeftPct     float64 `json:"swerve_left_pct" gorm:"column:left_pct"`
	RightPct    float64 `json:"swerve_right_pct" gorm:"column:right_pct"`
}

type TailgaterEffect struct {
	WithTailgater    *EffectMetric `json:"with_tailgater"`
	WithoutTailgater *EffectMetric `json:"without_tailgater"`
}

type EffectMetric struct {
	MaintainCount int64   `json:"maintain_count"`
	TotalCount    int64   `json:"total_count"`
	Percentage    float64 `json:"percentage"`
}

type ComplianceEffect struct {
	Compliant *EffectMetric `json:"compliant"`
	Violation *EffectMetric `json:"violation"`
}

type TimeDistributionPoint struct {
	Seconds int64 `json:"seconds" gorm:"column:seconds"`
	Count   int64 `json:"count" gorm:"column:count"`
}

type ArchetypeCount struct {
	Archetype string `json:"archetype" gorm:"column:archetype"`
	Count     int64  `json:"count" gorm:"column:count"`
}
