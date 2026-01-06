import type { ScenarioResponse } from '~/types/response.types'
import type { Entity } from '~/types/simulation'

type HighlightType = 'maintain' | 'swerve_left' | 'swerve_right'

export function useScenarioHighlight(
    scenario: Ref<ScenarioResponse | null>,
    highlightedZone: Ref<string | null>
) {
    const config = useRuntimeConfig()

    const tridentZoneDistance = computed(() => {
        return Number(config.public.tridentZoneDistance) || 1
    })

    const egoEntity = computed(() => 
        scenario.value?.entities.find(e => e.metadata.is_ego)
    )

    // Build approach path cells (between ego and zone start)
    const approachPathCells = computed(() => {
        if (!egoEntity.value) return new Set<string>()
        
        const ego = egoEntity.value
        const cells = new Set<string>()
        const distance = tridentZoneDistance.value
        
        let dRow = 0, dCol = 0
        switch (ego.metadata.orientation) {
            case 'N': dRow = -1; break
            case 'S': dRow = 1; break
            case 'E': dCol = 1; break
            case 'W': dCol = -1; break
        }
        
        for (let i = 1; i <= distance; i++) {
            const row = ego.row + (dRow * i)
            const col = ego.col + (dCol * i)
            cells.add(`${row},${col}`)
        }
        
        return cells
    })

    // Build lookup map: "row,col" -> highlight type
    const zoneCoordMap = computed(() => {
        const map = new Map<string, HighlightType>()
        if (!scenario.value?.trident_zones) return map
        
        const zones = scenario.value.trident_zones
        zones.zone_a.coordinates.forEach(c => map.set(`${c.row},${c.col}`, 'maintain'))
        zones.zone_b.coordinates.forEach(c => map.set(`${c.row},${c.col}`, 'swerve_left'))
        zones.zone_c.coordinates.forEach(c => map.set(`${c.row},${c.col}`, 'swerve_right'))
        
        return map
    })

    const zoneToType: Record<string, HighlightType> = {
        'zone_a': 'maintain',
        'zone_b': 'swerve_left',
        'zone_c': 'swerve_right'
    }

    function getHighlightType(row: number, col: number): HighlightType | null {
        if (!highlightedZone.value) return null
        
        const cellKey = `${row},${col}`
        const cellType = zoneCoordMap.value.get(cellKey)
        const activeType = zoneToType[highlightedZone.value]
        
        if (!activeType) return null
        if (cellType === activeType) return cellType
        if (approachPathCells.value.has(cellKey)) return activeType
        
        return null
    }

    function getEntitiesAt(row: number, col: number): Entity[] {
        return scenario.value?.entities.filter(e => e.row === row && e.col === col) || []
    }

    return {
        getHighlightType,
        getEntitiesAt,
        egoEntity
    }
}
