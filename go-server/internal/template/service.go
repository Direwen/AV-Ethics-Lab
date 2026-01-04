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
	CalculateTridentZones(templateID uuid.UUID, spawn domain.TridentSpawn) domain.TridentZones
	GetRearCoordinate(templateID uuid.UUID, row, col int, orientation domain.Direction) (*domain.EnrichedCoordinate, error)
}

type service struct {
	repo            Repository
	cache           []ContextTemplate
	laneConfigs     map[uuid.UUID]domain.LaneConfigMap
	tridentSpawns   map[uuid.UUID][]domain.TridentSpawn
	surfaceAt       map[uuid.UUID]map[[2]int]domain.SurfaceType // coord → surface
	laneDirectionAt map[uuid.UUID]map[[2]int]domain.Direction   // coord → lane direction
	gridDimensions  map[uuid.UUID][2]int                        // templateID → [height, width]
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
	gridDimensions := make(map[uuid.UUID][2]int)

	for _, template := range templates {
		var grid [][]int
		if err := json.Unmarshal(template.GridData, &grid); err != nil {
			return err
		}

		height := len(grid)
		width := len(grid[0])
		gridDimensions[template.Id] = [2]int{height, width}

		surfaceAt[template.Id] = make(map[[2]int]domain.SurfaceType)

		// Build surfaceAt lookup
		for row, cols := range grid {
			for col, tileID := range cols {
				if tile, exists := domain.TileRegistry[tileID]; exists {
					surfaceAt[template.Id][[2]int{row, col}] = tile.Definition.SurfaceType
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

		// Pre-compute valid trident spawns
		tridentSpawns[template.Id] = s.computeValidSpawns(
			template.Id,
			surfaceAt[template.Id],
			laneConfig,
			height,
			width,
		)
	}

	s.mu.Lock()
	s.cache = templates
	s.laneConfigs = laneConfigs
	s.tridentSpawns = tridentSpawns
	s.surfaceAt = surfaceAt
	s.laneDirectionAt = laneDirectionAt
	s.gridDimensions = gridDimensions
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

func (s *service) GetLaneDirectionAt(templateID uuid.UUID, row, col int) domain.Direction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if directions, ok := s.laneDirectionAt[templateID]; ok {
		return directions[[2]int{row, col}] // returns "" if not found
	}
	return ""
}

func (s *service) GetRearCoordinate(templateID uuid.UUID, row, col int, orientation domain.Direction) (*domain.EnrichedCoordinate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	surfaces, ok := s.surfaceAt[templateID]
	if !ok {
		return nil, errors.New("template not found")
	}

	directions := s.laneDirectionAt[templateID]
	dims := s.gridDimensions[templateID]
	height, width := dims[0], dims[1]

	// Rear is opposite of orientation
	var dRow, dCol int
	switch orientation {
	case domain.DirectionNorth:
		dRow, dCol = 1, 0
	case domain.DirectionSouth:
		dRow, dCol = -1, 0
	case domain.DirectionEast:
		dRow, dCol = 0, -1
	case domain.DirectionWest:
		dRow, dCol = 0, 1
	default:
		return nil, errors.New("invalid orientation")
	}

	rearRow := row + dRow
	rearCol := col + dCol

	if rearRow < 0 || rearRow >= height || rearCol < 0 || rearCol >= width {
		return nil, errors.New("rear coordinate out of bounds")
	}

	key := [2]int{rearRow, rearCol}
	surface := surfaces[key]

	if surface == domain.SurfaceDrivable {
		return &domain.EnrichedCoordinate{
			Coordinate:  domain.Coordinate{Row: rearRow, Col: rearCol},
			Surface:     surface,
			Orientation: directions[key],
		}, nil
	}

	return nil, errors.New("no valid rear coordinate found")
}

// CalculateTridentZones builds zones with expandable B/C (skips restricted, stops at building)
func (s *service) CalculateTridentZones(templateID uuid.UUID, spawn domain.TridentSpawn) domain.TridentZones {
	s.mu.RLock()
	defer s.mu.RUnlock()

	surfaces := s.surfaceAt[templateID]
	directions := s.laneDirectionAt[templateID]
	dims := s.gridDimensions[templateID]
	height, width := dims[0], dims[1]

	fRow, fCol, lRow, lCol, rRow, rCol := domain.CalculateTridentZones(spawn)

	const distance = 3 // how far ahead zones start
	const depth = 3    // how many traversable tiles to collect

	baseRow := spawn.Row + (fRow * distance)
	baseCol := spawn.Col + (fCol * distance)

	// Zone A: fixed strip forward (drivable + restricted allowed)
	generateForwardZone := func(startRow, startCol int) domain.TridentZone {
		coords := make([]domain.EnrichedCoordinate, 0, depth)
		for i := 0; i < depth; i++ {
			row := startRow + (fRow * i)
			col := startCol + (fCol * i)
			if row < 0 || row >= height || col < 0 || col >= width {
				break
			}
			key := [2]int{row, col}
			coords = append(coords, domain.EnrichedCoordinate{
				Coordinate:  domain.Coordinate{Row: row, Col: col},
				Surface:     surfaces[key],
				Orientation: directions[key],
			})
		}
		return domain.TridentZone{Coordinates: coords}
	}

	// Zone B/C: expand perpendicular, skip restricted, collect traversable, stop at building
	generateSideZone := func(perpRow, perpCol int) domain.TridentZone {
		coords := make([]domain.EnrichedCoordinate, 0, depth)

		// For each forward step, scan perpendicular until we find traversable
		for i := 0; i < depth; i++ {
			fwdRow := baseRow + (fRow * i)
			fwdCol := baseCol + (fCol * i)

			// Scan perpendicular from forward position
			for offset := 1; offset <= 5; offset++ { // max 5 tiles perpendicular
				row := fwdRow + (perpRow * offset)
				col := fwdCol + (perpCol * offset)

				if row < 0 || row >= height || col < 0 || col >= width {
					break // out of bounds
				}

				key := [2]int{row, col}
				surface := surfaces[key]

				if surface == domain.SurfaceBuilding {
					break // hit wall, stop scanning this direction
				}
				if surface == domain.SurfaceRestricted {
					continue // skip yellow lines
				}
				// Found traversable (drivable or walkable)
				coords = append(coords, domain.EnrichedCoordinate{
					Coordinate:  domain.Coordinate{Row: row, Col: col},
					Surface:     surface,
					Orientation: directions[key],
				})
				break // found one for this forward step, move to next
			}
		}
		return domain.TridentZone{Coordinates: coords}
	}

	return domain.TridentZones{
		ZoneA: generateForwardZone(baseRow, baseCol),
		ZoneB: generateSideZone(lRow, lCol),
		ZoneC: generateSideZone(rRow, rCol),
	}
}

// computeValidSpawns finds valid trident spawn points
func (s *service) computeValidSpawns(
	templateID uuid.UUID,
	surfaceAt map[[2]int]domain.SurfaceType,
	laneConfig domain.LaneConfigMap,
	height, width int,
) []domain.TridentSpawn {
	var validSpawns []domain.TridentSpawn

	for direction, coords := range laneConfig {
		for _, coord := range coords {
			// AV must spawn on drivable (not restricted)
			if surfaceAt[coord] != domain.SurfaceDrivable {
				continue
			}

			spawn := domain.TridentSpawn{
				Coordinate:  domain.Coordinate{Row: coord[0], Col: coord[1]},
				Orientation: direction,
			}

			if s.isValidSpawn(spawn, surfaceAt, height, width) {
				validSpawns = append(validSpawns, spawn)
			}
		}
	}

	return validSpawns
}

// isValidSpawn checks if spawn point produces valid trident zones
func (s *service) isValidSpawn(
	spawn domain.TridentSpawn,
	surfaceAt map[[2]int]domain.SurfaceType,
	height, width int,
) bool {
	fRow, fCol, lRow, lCol, rRow, rCol := domain.CalculateTridentZones(spawn)

	const distance = 3
	const depth = 3

	baseRow := spawn.Row + (fRow * distance)
	baseCol := spawn.Col + (fCol * distance)

	// Check Zone A: all cells must be road (drivable or restricted) and in bounds
	for i := 0; i < depth; i++ {
		row := baseRow + (fRow * i)
		col := baseCol + (fCol * i)
		if row < 0 || row >= height || col < 0 || col >= width {
			return false
		}
		surface := surfaceAt[[2]int{row, col}]
		if surface != domain.SurfaceDrivable && surface != domain.SurfaceRestricted {
			return false
		}
	}

	// Check Zone B and C: must find at least 1 traversable tile per zone
	checkSideZone := func(perpRow, perpCol int) bool {
		foundTraversable := 0
		for i := 0; i < depth; i++ {
			fwdRow := baseRow + (fRow * i)
			fwdCol := baseCol + (fCol * i)

			for offset := 1; offset <= 5; offset++ {
				row := fwdRow + (perpRow * offset)
				col := fwdCol + (perpCol * offset)

				if row < 0 || row >= height || col < 0 || col >= width {
					break
				}

				surface := surfaceAt[[2]int{row, col}]
				if surface == domain.SurfaceBuilding {
					break
				}
				if surface == domain.SurfaceRestricted {
					continue
				}
				// Found traversable
				foundTraversable++
				break
			}
		}
		return foundTraversable >= depth // need all 3 forward positions to have a side tile
	}

	if !checkSideZone(lRow, lCol) {
		return false
	}
	if !checkSideZone(rRow, rCol) {
		return false
	}

	return true
}
