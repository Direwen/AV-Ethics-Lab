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
	GetLaneConfig(templateID uuid.UUID) domain.LaneConfigMap
	GetRandomTridentSpawn(templateID uuid.UUID) (*domain.TridentSpawn, error)
	GetSurfaceAt(templateID uuid.UUID, row, col int) domain.SurfaceType
	GetLaneDirectionAt(templateID uuid.UUID, row, col int) domain.Direction
	EnrichTridentZones(templateID uuid.UUID, zones *domain.TridentZones)
}

type service struct {
	repo            Repository
	cache           []ContextTemplate
	laneConfigs     map[uuid.UUID]domain.LaneConfigMap
	tridentSpawns   map[uuid.UUID][]domain.TridentSpawn
	surfaceAt       map[uuid.UUID]map[[2]int]domain.SurfaceType // coord → surface
	laneDirectionAt map[uuid.UUID]map[[2]int]domain.Direction   // coord → lane direction
	mu              sync.RWMutex
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) LoadAllTemplates(ctx context.Context) error {
	templates, err := s.repo.GetAll(ctx)

	if err != nil {
		return err
	}

	laneConfigs := make(map[uuid.UUID]domain.LaneConfigMap)
	tridentSpawns := make(map[uuid.UUID][]domain.TridentSpawn)
	surfaceAt := make(map[uuid.UUID]map[[2]int]domain.SurfaceType)
	laneDirectionAt := make(map[uuid.UUID]map[[2]int]domain.Direction)

	for _, template := range templates {
		var grid [][]int
		if err := json.Unmarshal(template.GridData, &grid); err != nil {
			return err
		}

		surfaceAt[template.Id] = make(map[[2]int]domain.SurfaceType)

		// Build surfaceAt lookup and collect cells by surface for spawn validation
		cellsBySurface := make(map[domain.SurfaceType][][2]int)
		for row, cols := range grid {
			for col, tileID := range cols {
				if tile, exists := domain.TileRegistry[tileID]; exists {
					surface := tile.Definition.SurfaceType
					coord := [2]int{row, col}
					surfaceAt[template.Id][coord] = surface
					cellsBySurface[surface] = append(cellsBySurface[surface], coord)
				}
			}
		}

		// Parse lane config and build laneDirectionAt
		var laneConfig domain.LaneConfigMap
		if template.LaneConfig != nil {
			if err := json.Unmarshal(template.LaneConfig, &laneConfig); err != nil {
				return err
			}
		}
		laneConfigs[template.Id] = laneConfig

		laneDirectionAt[template.Id] = make(map[[2]int]domain.Direction)
		for direction, coords := range laneConfig {
			for _, coord := range coords {
				laneDirectionAt[template.Id][coord] = direction
			}
		}

		// Pre-compute valid trident spawns (uses local cellsBySurface)
		tridentSpawns[template.Id] = s.computeValidSpawns(
			cellsBySurface,
			laneConfig,
			len(grid),
			len(grid[0]),
		)
	}

	s.mu.Lock()
	s.cache = templates
	s.laneConfigs = laneConfigs
	s.tridentSpawns = tridentSpawns
	s.surfaceAt = surfaceAt
	s.laneDirectionAt = laneDirectionAt
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

func (s *service) GetLaneConfig(templateID uuid.UUID) domain.LaneConfigMap {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if config, ok := s.laneConfigs[templateID]; ok {
		return config
	}
	return nil
}

func (s *service) GetRandomTridentSpawn(templateID uuid.UUID) (*domain.TridentSpawn, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if spawns, ok := s.tridentSpawns[templateID]; ok {
		if len(spawns) == 0 {
			return nil, errors.New("no trident spawns found")
		}
		randomIndex := rand.Intn(len(spawns))
		return &spawns[randomIndex], nil
	}
	return nil, errors.New("no trident spawns found")
}

// GetSurfaceAt returns the surface type at a coordinate (O(1) lookup)
func (s *service) GetSurfaceAt(templateID uuid.UUID, row, col int) domain.SurfaceType {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if surfaces, ok := s.surfaceAt[templateID]; ok {
		if surface, exists := surfaces[[2]int{row, col}]; exists {
			return surface
		}
	}
	return domain.SurfaceBuilding // default
}

// GetLaneDirectionAt returns the lane direction at a coordinate (O(1) lookup)
func (s *service) GetLaneDirectionAt(templateID uuid.UUID, row, col int) domain.Direction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if directions, ok := s.laneDirectionAt[templateID]; ok {
		return directions[[2]int{row, col}] // returns "" if not found
	}
	return ""
}

// EnrichTridentZones fills in Surface and Orientation for each coordinate in zones
func (s *service) EnrichTridentZones(templateID uuid.UUID, zones *domain.TridentZones) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	surfaces := s.surfaceAt[templateID]
	directions := s.laneDirectionAt[templateID]

	enrichZone := func(zone *domain.TridentZone) {
		for i := range zone.Coordinates {
			coord := &zone.Coordinates[i]
			key := [2]int{coord.Row, coord.Col}
			coord.Surface = surfaces[key]
			coord.Orientation = directions[key]
		}
	}

	enrichZone(&zones.ZoneA)
	enrichZone(&zones.ZoneB)
	enrichZone(&zones.ZoneC)
}

// computeValidSpawns finds valid trident spawn points using pre-indexed surface data
func (s *service) computeValidSpawns(
	cellsBySurface map[domain.SurfaceType][][2]int,
	laneConfig domain.LaneConfigMap,
	height, width int,
) []domain.TridentSpawn {
	var validSpawns []domain.TridentSpawn

	// Build lookup sets for O(1) surface checks
	drivableSet := s.coordSet(cellsBySurface[domain.SurfaceDrivable])
	walkableSet := s.coordSet(cellsBySurface[domain.SurfaceWalkable])
	restrictedSet := s.coordSet(cellsBySurface[domain.SurfaceRestricted])

	// Road = drivable + restricted (yellow lines are still road)
	roadSet := s.mergeCoordSets(drivableSet, restrictedSet)
	// Traversable = anything that's not a building
	traversableSet := s.mergeCoordSets(drivableSet, walkableSet, restrictedSet)

	// Only check lane cells — direction comes from lane config
	for direction, coords := range laneConfig {
		for _, coord := range coords {
			spawn := domain.TridentSpawn{
				Coordinate:  domain.Coordinate{Row: coord[0], Col: coord[1]},
				Orientation: direction,
			}

			zones := domain.CalculateTridentZones(spawn)

			if s.IsValidTrident(zones, roadSet, traversableSet, height, width) {
				validSpawns = append(validSpawns, spawn)
			}
		}
	}

	return validSpawns
}

// isValidTrident checks if all three zones land on valid terrain
func (s *service) IsValidTrident(
	zones domain.TridentZones,
	roadSet, traversableSet map[[2]int]bool,
	height, width int,
) bool {
	// Zone A (forward): must be road
	for _, coord := range zones.ZoneA.Coordinates {
		if !s.inBounds(coord.Coordinate, height, width) || !roadSet[[2]int{coord.Row, coord.Col}] {
			return false
		}
	}

	// Zone B (left): must be traversable
	for _, coord := range zones.ZoneB.Coordinates {
		if !s.inBounds(coord.Coordinate, height, width) || !traversableSet[[2]int{coord.Row, coord.Col}] {
			return false
		}
	}

	// Zone C (right): must be traversable
	for _, coord := range zones.ZoneC.Coordinates {
		if !s.inBounds(coord.Coordinate, height, width) || !traversableSet[[2]int{coord.Row, coord.Col}] {
			return false
		}
	}

	return true
}

// coordSet converts a slice of coordinates to a map for O(1) lookup
func (s *service) coordSet(coords [][2]int) map[[2]int]bool {
	set := make(map[[2]int]bool, len(coords))
	for _, c := range coords {
		set[c] = true
	}
	return set
}

// mergeCoordSets combines multiple coordinate sets into one
func (s *service) mergeCoordSets(sets ...map[[2]int]bool) map[[2]int]bool {
	total := 0
	for _, set := range sets {
		total += len(set)
	}
	merged := make(map[[2]int]bool, total)
	for _, set := range sets {
		for k := range set {
			merged[k] = true
		}
	}
	return merged
}

func (s *service) inBounds(coord domain.Coordinate, height, width int) bool {
	return coord.Row >= 0 && coord.Row < height && coord.Col >= 0 && coord.Col < width
}
