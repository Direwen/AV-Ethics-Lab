package session

import (
	"context"

	"github.com/direwen/go-server/internal/util"
)

type Service interface {
	RegisterSession(ctx context.Context, input CreateSessionInput) (string, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RegisterSession(ctx context.Context, input CreateSessionInput) (string, error) {
	exists, err := s.repo.FingerprintExists(ctx, input.Fingerprint)
	if err != nil {
		return "", err
	}

	sess := Session{
		AgeRange:          input.AgeRange,
		Gender:            input.Gender,
		Country:           input.Country,
		Occupation:        input.Occupation,
		DrivingExperience: input.DrivingExperience,
		Fingerprint:       input.Fingerprint,
		SelfReportedNew:   input.SelfReportedNew,
		IsDuplicate:       exists,
	}

	if err := s.repo.Create(ctx, &sess); err != nil {
		return "", err
	}

	signedToken, err := util.GenerateToken(
		sess.Id.String(),
		map[string]any{"issuer": "av-ethics-lab"},
	)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
