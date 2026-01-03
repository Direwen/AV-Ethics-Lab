package template

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"sync"

	"github.com/direwen/go-server/internal/shared/domain"
	"github.com/google/uuid"
)

type Service interface {
	LoadAllTemplates(ctx context.Context) error
	GetAllTemplates(ctx context.Context) ([]ContextTemplate, error)
	GetByID(id uuid.UUID) (*ContextTemplate, error)
	PickTemplate(excludeIDs []uuid.UUID) (*ContextTemplate, error)
	GetCellsBySurface(templateID uuid.UUID, surface domain.SurfaceType) [][2]int
	GetLaneConfig(templateID uuid.UUID) domain.LaneConfigMap
}

type service struct {
	repo           Repository
	cache          []ContextTemplate
	cellsBySurface map[uuid.UUID]map[domain.SurfaceType][][2]int
	laneConfigs    map[uuid.UUID]domain.LaneConfigMap
	mu             sync.RWMutex //only one write | multiple reads
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) LoadAllTemplates(ctx context.Context) error {
	templates, err := s.repo.GetAll(ctx)

	if err != nil {
		return err
	}

	cellsBySurface := make(map[uuid.UUID]map[domain.SurfaceType][][2]int)
	laneConfigs := make(map[uuid.UUID]domain.LaneConfigMap)

	for _, template := range templates {
		var grid [][]int
		if err := json.Unmarshal(template.GridData, &grid); err != nil {
			return err
		}
		cellsBySurface[template.Id] = make(map[domain.SurfaceType][][2]int)
		for row, cols := range grid {
			for col, tileID := range cols {
				if tile, exists := domain.TileRegistry[tileID]; exists {
					surface := tile.Definition.SurfaceType
					cellsBySurface[template.Id][surface] = append(
						cellsBySurface[template.Id][surface],
						[2]int{row, col},
					)
				}
			}
		}

		// Parse lane config
		var laneConfig domain.LaneConfigMap
		if template.LaneConfig != nil {
			if err := json.Unmarshal(template.LaneConfig, &laneConfig); err != nil {
				return err
			}
		}
		laneConfigs[template.Id] = laneConfig
	}

	s.mu.Lock()
	s.cache = templates
	s.cellsBySurface = cellsBySurface
	s.laneConfigs = laneConfigs
	s.mu.Unlock()

	return nil
}

func (s *service) GetAllTemplates(ctx context.Context) ([]ContextTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.cache) == 0 {
		return nil, errors.New("no templates found")
	}

	return s.cache, nil
}

func (s *service) GetByID(id uuid.UUID) (*ContextTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, t := range s.cache {
		if t.Id == id {
			return &t, nil
		}
	}
	return nil, errors.New("template not found")
}

func (s *service) PickTemplate(excludeIDs []uuid.UUID) (*ContextTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.cache) == 0 {
		return nil, errors.New("template cache is empty")
	}

	var candidates []ContextTemplate
	excludeIDsMap := make(map[uuid.UUID]bool, len(excludeIDs))
	for _, id := range excludeIDs {
		excludeIDsMap[id] = true
	}

	for _, t := range s.cache {
		if _, found := excludeIDsMap[t.Id]; !found {
			candidates = append(candidates, t)
		}
	}

	if len(candidates) == 0 {
		candidates = s.cache
	}

	randomIndex := rand.Intn(len(candidates))
	return &candidates[randomIndex], nil
}

func (s *service) GetCellsBySurface(templateID uuid.UUID, surface domain.SurfaceType) [][2]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if cells, ok := s.cellsBySurface[templateID]; ok {
		return cells[surface]
	}
	return nil
}

func (s *service) GetLaneConfig(templateID uuid.UUID) domain.LaneConfigMap {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if config, ok := s.laneConfigs[templateID]; ok {
		return config
	}
	return nil
}
