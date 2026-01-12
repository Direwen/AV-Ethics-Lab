package dashboard

import (
	"context"
	"strings"

	"github.com/direwen/go-server/internal/response"
	"github.com/direwen/go-server/internal/scenario"
	"github.com/direwen/go-server/internal/session"
	"github.com/direwen/go-server/internal/shared/models"
	"github.com/direwen/go-server/pkg/database"
)

type Service interface {
	GetPublicStats(ctx context.Context) (*PublicStats, error)
}

type service struct {
	sessionRepo  session.Repository
	scenarioRepo scenario.Repository
	responseRepo response.Repository
}

func NewService(sessionRepo session.Repository, scenarioRepo scenario.Repository, responseRepo response.Repository) Service {
	return &service{
		sessionRepo:  sessionRepo,
		scenarioRepo: scenarioRepo,
		responseRepo: responseRepo,
	}
}

func (s *service) GetPublicStats(ctx context.Context) (*PublicStats, error) {
	sessionCount, err := s.sessionRepo.CountSessions(ctx, database.WithFilter("sessions.status = ?", models.StatusCompleted))
	if err != nil {
		return nil, err
	}

	countriesCount, err := s.sessionRepo.CountCountries(ctx)
	if err != nil {
		return nil, err
	}

	leastHarmfulRaw, err := s.responseRepo.CountRankingAtIndex(ctx, 0)
	if err != nil {
		return nil, err
	}
	leastHarmfulOutcome := toOutcomeDistribution(leastHarmfulRaw)

	selfPreservationRaw, err := s.responseRepo.GetCountsByFactor(ctx, "has_tailgater", 0)
	if err != nil {
		return nil, err
	}
	selfPreservationEffect := toTailgaterEffect(selfPreservationRaw)

	entityComplianceRaw, err := s.responseRepo.GetCountsByFactor(ctx, "primary_behavior", 0)
	if err != nil {
		return nil, err
	}
	entityComplianceEffect := toComplianceEffect(entityComplianceRaw)

	responseTimeRaw, err := s.responseRepo.GetResponseTimeDistribution(ctx)
	if err != nil {
		return nil, err
	}
	decisionTimeDistribution := toTimeDistribution(responseTimeRaw)

	archetypeRaw, err := s.sessionRepo.CountArchetypes(ctx)
	if err != nil {
		return nil, err
	}
	archetypeDistribution := toArchetypeDistribution(archetypeRaw)

	return &PublicStats{
		CompletedSessions:        sessionCount,
		CountriesRepresented:     countriesCount,
		LeastHarmfulOutcome:      leastHarmfulOutcome,
		SelfPreservationEffect:   selfPreservationEffect,
		EntityComplianceEffect:   entityComplianceEffect,
		DecisionTimeDistribution: decisionTimeDistribution,
		ArchetypeDistribution:    archetypeDistribution,
	}, nil
}

func toOutcomeDistribution(raw map[string]int64) *OutcomeDistribution {
	maintain := raw["maintain"]
	swerveLeft := raw["swerve_left"]
	swerveRight := raw["swerve_right"]
	total := maintain + swerveLeft + swerveRight

	var maintainPct, leftPct, rightPct float64
	if total > 0 {
		maintainPct = float64(maintain) / float64(total) * 100
		leftPct = float64(swerveLeft) / float64(total) * 100
		rightPct = float64(swerveRight) / float64(total) * 100
	}

	return &OutcomeDistribution{
		Maintain:    maintain,
		SwerveLeft:  swerveLeft,
		SwerveRight: swerveRight,
		Total:       total,
		MaintainPct: maintainPct,
		LeftPct:     leftPct,
		RightPct:    rightPct,
	}
}

func toTailgaterEffect(raw []response.FactorCount) *TailgaterEffect {
	withTailgater := &EffectMetric{}
	withoutTailgater := &EffectMetric{}

	for _, fc := range raw {
		// Normalize factor value to handle JSON variations
		normalized := strings.ToLower(strings.TrimSpace(fc.FactorValue))
		hasTailgater := normalized == "true" || normalized == "1" || normalized == "yes"

		if hasTailgater {
			withTailgater.TotalCount += fc.Count
			if fc.Outcome == "maintain" {
				withTailgater.MaintainCount += fc.Count
			}
		} else {
			withoutTailgater.TotalCount += fc.Count
			if fc.Outcome == "maintain" {
				withoutTailgater.MaintainCount += fc.Count
			}
		}
	}

	if withTailgater.TotalCount > 0 {
		withTailgater.Percentage = float64(withTailgater.MaintainCount) / float64(withTailgater.TotalCount) * 100
	}
	if withoutTailgater.TotalCount > 0 {
		withoutTailgater.Percentage = float64(withoutTailgater.MaintainCount) / float64(withoutTailgater.TotalCount) * 100
	}

	return &TailgaterEffect{
		WithTailgater:    withTailgater,
		WithoutTailgater: withoutTailgater,
	}
}

func toComplianceEffect(raw []response.FactorCount) *ComplianceEffect {
	compliant := &EffectMetric{}
	violation := &EffectMetric{}

	for _, fc := range raw {
		normalized := strings.ToLower(strings.TrimSpace(fc.FactorValue))
		isCompliant := normalized == "compliant"
		isViolation := normalized == "violation"

		if isCompliant {
			compliant.TotalCount += fc.Count
			if fc.Outcome == "maintain" {
				compliant.MaintainCount += fc.Count
			}
		} else if isViolation {
			violation.TotalCount += fc.Count
			if fc.Outcome == "maintain" {
				violation.MaintainCount += fc.Count
			}
		}
	}

	if compliant.TotalCount > 0 {
		compliant.Percentage = float64(compliant.MaintainCount) / float64(compliant.TotalCount) * 100
	}
	if violation.TotalCount > 0 {
		violation.Percentage = float64(violation.MaintainCount) / float64(violation.TotalCount) * 100
	}

	return &ComplianceEffect{
		Compliant: compliant,
		Violation: violation,
	}
}

func toTimeDistribution(raw []response.ResponseTimeMSCount) *TimeDistribution {
	dist := &TimeDistribution{}

	for _, r := range raw {
		switch {
		case r.Second < 2:
			dist.Under2Sec += r.Count
		case r.Second < 4:
			dist.Between2And4 += r.Count
		case r.Second < 6:
			dist.Between4And6 += r.Count
		default:
			dist.Over6Sec += r.Count
		}
		dist.Total += r.Count
	}

	return dist
}

func toArchetypeDistribution(raw []session.ArchetypeCount) []ArchetypeCount {
	result := make([]ArchetypeCount, len(raw))
	for i, a := range raw {
		result[i] = ArchetypeCount{
			Archetype: a.Archetype,
			Count:     a.Count,
		}
	}
	return result
}
