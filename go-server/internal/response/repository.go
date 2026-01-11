package response

import (
	"context"
	"fmt"

	"github.com/direwen/go-server/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, response *Response) error
	GetByID(ctx context.Context, id uuid.UUID, opts ...database.QueryOption) (*Response, error)
	GetByScenarioID(ctx context.Context, scenarioID uuid.UUID, opts ...database.QueryOption) (*Response, error)
	CountBySessionID(ctx context.Context, sessionID uuid.UUID) (int, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID, opts ...database.QueryOption) ([]*Response, error)
	CountRankingAtIndex(ctx context.Context, index int) (map[string]int64, error)
	GetCountsByFactor(ctx context.Context, factorKey string, rankIndex int) ([]FactorCount, error)
	GetResponseTimeDistribution(ctx context.Context) ([]ResponseTimeMSCount, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, response *Response) error {
	return database.GetDB(ctx, r.db).WithContext(ctx).Create(response).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID, opts ...database.QueryOption) (*Response, error) {
	var response Response
	db := database.GetDB(ctx, r.db).WithContext(ctx).Model(&Response{}).Where("id = ?", id)
	db = database.ApplyOptions(db, opts...)
	err := db.First(&response).Error
	return &response, err
}

func (r *repository) GetByScenarioID(ctx context.Context, scenarioID uuid.UUID, opts ...database.QueryOption) (*Response, error) {
	var response Response
	db := database.GetDB(ctx, r.db).WithContext(ctx).Model(&Response{}).Where("scenario_id = ?", scenarioID)
	db = database.ApplyOptions(db, opts...)
	err := db.First(&response).Error
	return &response, err
}

func (r *repository) CountBySessionID(ctx context.Context, sessionID uuid.UUID) (int, error) {
	var count int64
	err := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&Response{}).
		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
		Where("scenarios.session_id = ?", sessionID).
		Count(&count).Error
	return int(count), err
}

func (r *repository) GetBySessionID(ctx context.Context, sessionID uuid.UUID, opts ...database.QueryOption) ([]*Response, error) {
	var responses []*Response
	db := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&Response{}).
		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
		Where("scenarios.session_id = ?", sessionID)
	db = database.ApplyOptions(db, opts...)
	err := db.Find(&responses).Error
	return responses, err
}

func (r *repository) CountRankingAtIndex(ctx context.Context, index int) (map[string]int64, error) {

	var results []RankResult

	jsonSelector := fmt.Sprintf("ranking_order->>%d", index)
	err := database.GetDB(ctx, r.db).
		WithContext(ctx).
		Model(&Response{}).
		Select(jsonSelector+" as outcome", "COUNT(*) as count").
		Group(jsonSelector).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, res := range results {
		counts[res.Outcome] = res.Count
	}

	return counts, nil
}

func (r *repository) GetCountsByFactor(ctx context.Context, factorKey string, rankIndex int) ([]FactorCount, error) {
	var results []FactorCount

	factorSelector := fmt.Sprintf("scenarios.factors->>'%s'", factorKey)
	rankSelector := fmt.Sprintf("responses.ranking_order->>%d", rankIndex)

	err := database.GetDB(ctx, r.db).
		WithContext(ctx).
		Table("responses").
		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
		Joins("JOIN sessions ON sessions.id = scenarios.session_id").
		Where("sessions.status = ?", "completed").
		Select(fmt.Sprintf("%s as factor_value, %s as outcome, COUNT(*) as total", factorSelector, rankSelector)).
		Group("factor_value, outcome").
		Scan(&results).Error

	return results, err
}

func (r *repository) GetResponseTimeDistribution(ctx context.Context) ([]ResponseTimeMSCount, error) {

	var results []ResponseTimeMSCount

	err := database.GetDB(ctx, r.db).
		WithContext(ctx).
		Model(&Response{}).
		Where("is_timeout = ?", false).
		Select("FLOOR(response_time_ms/1000) as second, COUNT(*) as count").
		Group("second").
		Order("second ASC").
		Scan(&results).Error

	return results, err
}

// ========== INEFFICIENT QUERY =========
// func (r *repository) GetRankingByFactor(
// 	ctx context.Context,
// 	factorKey string,
// 	factorValue string,
// 	rankIndex int,
// 	targetOutcome string,
// ) (int64, error) {

// 	jsonFactorSelector := fmt.Sprintf("scenarios.factor->>%s", factorKey)
// 	jsonRankingSelector := fmt.Sprintf("responses.ranking_order->>%d", rankIndex)
// 	var count int64

// 	err := database.GetDB(ctx, r.db).
// 		WithContext(ctx).
// 		Model(&Response{}).
// 		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
// 		Joins("JOIN sessions ON sessions.id = scenarios.session_id").
// 		Where("sessions.status = ?", "completed").
// 		Where(jsonFactorSelector+" = ?", factorValue).
// 		Where(jsonRankingSelector+" = ?", targetOutcome).
// 		Count(&count).Error

// 	return count, err
// }
