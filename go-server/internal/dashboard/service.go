package dashboard

import (
	"context"
)

type Service interface {
	GetPublicStats(ctx context.Context) (*PublicStats, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetPublicStats(ctx context.Context) (*PublicStats, error) {
	sessionCount, err := s.repo.GetCompletedSessionCount(ctx)
	if err != nil {
		return nil, err
	}

	countriesCount, err := s.repo.GetCountryCount(ctx)
	if err != nil {
		return nil, err
	}

	leastHarmfulOutcome, err := s.repo.GetLeastHarmfulOutcome(ctx)
	if err != nil {
		return nil, err
	}

	selfPreservationEffect, err := s.repo.GetTailgaterEffect(ctx)
	if err != nil {
		return nil, err
	}

	entityComplianceEffect, err := s.repo.GetComplianceEffect(ctx)
	if err != nil {
		return nil, err
	}

	decisionTimeDistribution, err := s.repo.GetTimeDistribution(ctx)
	if err != nil {
		return nil, err
	}

	archetypeDistribution, err := s.repo.GetArchetypeDistribution(ctx)
	if err != nil {
		return nil, err
	}

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
