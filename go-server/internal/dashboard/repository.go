package dashboard

import (
	"context"
	"fmt"

	"github.com/direwen/go-server/internal/shared/models"
	"github.com/direwen/go-server/pkg/database"
	"gorm.io/gorm"
)

// JSON path selectors
const (
	jsonRankingFirst    = "responses.ranking_order->>0"
	jsonHasTailgater    = "scenarios.factors->>'has_tailgater'"
	jsonPrimaryBehavior = "scenarios.factors->>'primary_behavior'"
	jsonArchetype       = "feedback->>'archetype'"
)

// Action values
const (
	actionMaintain    = "maintain"
	actionSwerveLeft  = "swerve_left"
	actionSwerveRight = "swerve_right"
)

// Factor values
const (
	factorTrue      = "true"
	factorFalse     = "false"
	factorCompliant = "compliant"
	factorViolation = "violation"
)

type Repository interface {
	GetCompletedSessionCount(ctx context.Context) (int64, error)
	GetCountryCount(ctx context.Context) (int64, error)
	GetLeastHarmfulOutcome(ctx context.Context) (*OutcomeDistribution, error)
	GetTailgaterEffect(ctx context.Context) (*TailgaterEffect, error)
	GetComplianceEffect(ctx context.Context) (*ComplianceEffect, error)
	GetTimeDistribution(ctx context.Context) ([]TimeDistributionPoint, error)
	GetArchetypeDistribution(ctx context.Context) ([]ArchetypeCount, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetCompletedSessionCount(ctx context.Context) (int64, error) {
	var count int64
	err := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&models.Session{}).
		Where("status = ?", models.StatusCompleted).
		Count(&count).Error
	return count, err
}

func (r *repository) GetCountryCount(ctx context.Context) (int64, error) {
	var count int64
	err := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&models.Session{}).
		Where("status = ?", models.StatusCompleted).
		Distinct("country").
		Count(&count).Error
	return count, err
}

func (r *repository) GetLeastHarmfulOutcome(ctx context.Context) (*OutcomeDistribution, error) {
	var result OutcomeDistribution

	query := fmt.Sprintf(`
		COUNT(*) FILTER (WHERE %s = '%s') as maintain,
		COUNT(*) FILTER (WHERE %s = '%s') as swerve_left,
		COUNT(*) FILTER (WHERE %s = '%s') as swerve_right,
		COUNT(*) as total,
		COALESCE(100.0 * COUNT(*) FILTER (WHERE %s = '%s') / NULLIF(COUNT(*), 0), 0) as maintain_pct,
		COALESCE(100.0 * COUNT(*) FILTER (WHERE %s = '%s') / NULLIF(COUNT(*), 0), 0) as left_pct,
		COALESCE(100.0 * COUNT(*) FILTER (WHERE %s = '%s') / NULLIF(COUNT(*), 0), 0) as right_pct
	`,
		jsonRankingFirst, actionMaintain,
		jsonRankingFirst, actionSwerveLeft,
		jsonRankingFirst, actionSwerveRight,
		jsonRankingFirst, actionMaintain,
		jsonRankingFirst, actionSwerveLeft,
		jsonRankingFirst, actionSwerveRight,
	)

	err := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&models.Response{}).
		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
		Joins("JOIN sessions ON sessions.id = scenarios.session_id").
		Where("sessions.status = ?", models.StatusCompleted).
		Where("responses.has_interacted = ?", true).
		Select(query).
		Scan(&result).Error

	return &result, err
}

func (r *repository) GetTailgaterEffect(ctx context.Context) (*TailgaterEffect, error) {
	type rawResult struct {
		WithMaintain    int64   `gorm:"column:with_maintain"`
		WithTotal       int64   `gorm:"column:with_total"`
		WithPct         float64 `gorm:"column:with_pct"`
		WithoutMaintain int64   `gorm:"column:without_maintain"`
		WithoutTotal    int64   `gorm:"column:without_total"`
		WithoutPct      float64 `gorm:"column:without_pct"`
	}

	query := fmt.Sprintf(`
		COUNT(*) FILTER (WHERE LOWER(%s) = '%s' AND %s = '%s') as with_maintain,
		COUNT(*) FILTER (WHERE LOWER(%s) = '%s') as with_total,
		COALESCE(100.0 * COUNT(*) FILTER (WHERE LOWER(%s) = '%s' AND %s = '%s') / NULLIF(COUNT(*) FILTER (WHERE LOWER(%s) = '%s'), 0), 0) as with_pct,
		COUNT(*) FILTER (WHERE LOWER(%s) = '%s' AND %s = '%s') as without_maintain,
		COUNT(*) FILTER (WHERE LOWER(%s) = '%s') as without_total,
		COALESCE(100.0 * COUNT(*) FILTER (WHERE LOWER(%s) = '%s' AND %s = '%s') / NULLIF(COUNT(*) FILTER (WHERE LOWER(%s) = '%s'), 0), 0) as without_pct
	`,
		jsonHasTailgater, factorTrue, jsonRankingFirst, actionMaintain,
		jsonHasTailgater, factorTrue,
		jsonHasTailgater, factorTrue, jsonRankingFirst, actionMaintain, jsonHasTailgater, factorTrue,
		jsonHasTailgater, factorFalse, jsonRankingFirst, actionMaintain,
		jsonHasTailgater, factorFalse,
		jsonHasTailgater, factorFalse, jsonRankingFirst, actionMaintain, jsonHasTailgater, factorFalse,
	)

	var raw rawResult
	err := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&models.Response{}).
		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
		Joins("JOIN sessions ON sessions.id = scenarios.session_id").
		Where("sessions.status = ?", models.StatusCompleted).
		Where("responses.has_interacted = ?", true).
		Select(query).
		Scan(&raw).Error

	if err != nil {
		return nil, err
	}

	return &TailgaterEffect{
		WithTailgater: &EffectMetric{
			MaintainCount: raw.WithMaintain,
			TotalCount:    raw.WithTotal,
			Percentage:    raw.WithPct,
		},
		WithoutTailgater: &EffectMetric{
			MaintainCount: raw.WithoutMaintain,
			TotalCount:    raw.WithoutTotal,
			Percentage:    raw.WithoutPct,
		},
	}, nil
}

func (r *repository) GetComplianceEffect(ctx context.Context) (*ComplianceEffect, error) {
	type rawResult struct {
		CompliantMaintain int64   `gorm:"column:compliant_maintain"`
		CompliantTotal    int64   `gorm:"column:compliant_total"`
		CompliantPct      float64 `gorm:"column:compliant_pct"`
		ViolationMaintain int64   `gorm:"column:violation_maintain"`
		ViolationTotal    int64   `gorm:"column:violation_total"`
		ViolationPct      float64 `gorm:"column:violation_pct"`
	}

	query := fmt.Sprintf(`
		COUNT(*) FILTER (WHERE LOWER(%s) = '%s' AND %s = '%s') as compliant_maintain,
		COUNT(*) FILTER (WHERE LOWER(%s) = '%s') as compliant_total,
		COALESCE(100.0 * COUNT(*) FILTER (WHERE LOWER(%s) = '%s' AND %s = '%s') / NULLIF(COUNT(*) FILTER (WHERE LOWER(%s) = '%s'), 0), 0) as compliant_pct,
		COUNT(*) FILTER (WHERE LOWER(%s) = '%s' AND %s = '%s') as violation_maintain,
		COUNT(*) FILTER (WHERE LOWER(%s) = '%s') as violation_total,
		COALESCE(100.0 * COUNT(*) FILTER (WHERE LOWER(%s) = '%s' AND %s = '%s') / NULLIF(COUNT(*) FILTER (WHERE LOWER(%s) = '%s'), 0), 0) as violation_pct
	`,
		jsonPrimaryBehavior, factorCompliant, jsonRankingFirst, actionMaintain,
		jsonPrimaryBehavior, factorCompliant,
		jsonPrimaryBehavior, factorCompliant, jsonRankingFirst, actionMaintain, jsonPrimaryBehavior, factorCompliant,
		jsonPrimaryBehavior, factorViolation, jsonRankingFirst, actionMaintain,
		jsonPrimaryBehavior, factorViolation,
		jsonPrimaryBehavior, factorViolation, jsonRankingFirst, actionMaintain, jsonPrimaryBehavior, factorViolation,
	)

	var raw rawResult
	err := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&models.Response{}).
		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
		Joins("JOIN sessions ON sessions.id = scenarios.session_id").
		Where("sessions.status = ?", models.StatusCompleted).
		Where("responses.has_interacted = ?", true).
		Select(query).
		Scan(&raw).Error

	if err != nil {
		return nil, err
	}

	return &ComplianceEffect{
		Compliant: &EffectMetric{
			MaintainCount: raw.CompliantMaintain,
			TotalCount:    raw.CompliantTotal,
			Percentage:    raw.CompliantPct,
		},
		Violation: &EffectMetric{
			MaintainCount: raw.ViolationMaintain,
			TotalCount:    raw.ViolationTotal,
			Percentage:    raw.ViolationPct,
		},
	}, nil
}

func (r *repository) GetTimeDistribution(ctx context.Context) ([]TimeDistributionPoint, error) {
	var result []TimeDistributionPoint

	err := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&models.Response{}).
		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
		Joins("JOIN sessions ON sessions.id = scenarios.session_id").
		Where("sessions.status = ?", models.StatusCompleted).
		Where("responses.has_interacted = ?", true).
		Where("responses.is_timeout = ?", false).
		Select(`
			FLOOR(responses.response_time_ms / 1000) as seconds,
			COUNT(*) as count
		`).
		Group("seconds").
		Scan(&result).Error

	return result, err
}

func (r *repository) GetArchetypeDistribution(ctx context.Context) ([]ArchetypeCount, error) {
	var results []ArchetypeCount

	err := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&models.Session{}).
		Where("status = ?", models.StatusCompleted).
		Where(jsonArchetype + " IS NOT NULL").
		Select(jsonArchetype + " as archetype, COUNT(*) as count").
		Group(jsonArchetype).
		Order("count DESC").
		Scan(&results).Error

	return results, err
}
