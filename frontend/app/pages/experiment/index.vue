<template>
    <div class="container mx-auto px-4 py-6">
        <!-- Header with Timer -->
        <div class="flex items-center justify-between mb-6">
            <div>
                <h1 class="text-xl font-semibold">Scenario {{ scenario.current_step }} of {{ scenario.total_steps }}</h1>
                <p class="text-sm text-[hsl(var(--maz-muted))]">{{ scenario.template_name }}</p>
            </div>
            <ExperimentTimer :duration="30" loop />
        </div>

        <!-- Narrative -->
        <div class="mb-6 p-4 rounded-xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-warning))]/10">
            <p class="text-sm leading-relaxed">{{ scenario.narrative }}</p>
        </div>

        <!-- Main Content -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 lg:gap-8">
            <!-- Left: Scenario Board (2/3) -->
            <section class="flex flex-col lg:col-span-2">
                <h2 class="text-sm font-medium uppercase tracking-wider text-[hsl(var(--maz-muted))] mb-3">
                    Scenario
                </h2>
                <SimulationScenarioContainer :visibility="scenario.factors.visibility as 'Clear' | 'Fog' | 'Night' | 'Rain'">
                    <div v-for="(row, rIndex) in scenario.grid_data" :key="rIndex" class="flex">
                        <SimulationCell
                            v-for="(cellCode, cIndex) in row"
                            :key="cIndex"
                            :cell-code="cellCode"
                            :definition="getCellDefinition(String(cellCode))"
                            :entities="getEntitiesAt(rIndex, cIndex)"
                            :highlight-type="getHighlightType(rIndex, cIndex)"
                            :road-condition="scenario.factors.road_condition as 'Dry' | 'Wet' | 'Icy'"
                        />
                    </div>
                </SimulationScenarioContainer>

                <div class="flex flex-wrap gap-2 mt-3">
                    <span 
                        class="px-2 py-1 rounded-full text-xs font-medium"
                        :class="speedBadgeClass"
                    >
                        {{ getSpeedIcon(scenario.factors.speed) }} {{ scenario.factors.speed }} Speed
                    </span>
                    <span 
                        class="px-2 py-1 rounded-full text-xs font-medium"
                        :class="brakeBadgeClass"
                    >
                        {{ getBrakeIcon(scenario.factors.brake_status) }} Brake: {{ scenario.factors.brake_status }}
                    </span>
                    <span 
                        v-if="scenario.factors.has_tailgater"
                        class="px-2 py-1 rounded-full text-xs font-medium bg-[hsl(var(--maz-danger))]/20 text-[hsl(var(--maz-danger))]"
                    >
                        🚗 Tailgater
                    </span>
                </div>
            </section>

            <!-- Right: Ranking Options -->
            <section class="flex flex-col">
                <h2 class="text-sm font-medium uppercase tracking-wider text-[hsl(var(--maz-muted))] mb-3">
                    Rank Options
                </h2>
                <div class="flex-1 flex flex-col gap-3">
                    <div 
                        v-for="(option, index) in rankingOptions" 
                        :key="option.key"
                        class="flex items-center gap-4 p-4 rounded-xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-secondary))] cursor-grab hover:border-[hsl(var(--maz-primary))] transition-colors"
                        @mouseenter="highlightedZone = option.zone"
                        @mouseleave="highlightedZone = null"
                    >
                        <span class="flex items-center justify-center w-8 h-8 rounded-full bg-[hsl(var(--maz-primary))]/20 text-[hsl(var(--maz-primary))] font-bold text-sm">
                            {{ index + 1 }}
                        </span>
                        <div class="flex-1">
                            <p class="font-medium">{{ option.label }}</p>
                        </div>
                        <MazIcon name="bars-3" class="w-5 h-5 text-[hsl(var(--maz-muted))]" />
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useCellDefinition } from '~/composables/useCellDefinition'
import type { Entity } from '~/types/simulation'

definePageMeta({
    layout: 'optional-layout',
    middleware: ['session']
})

const config = useRuntimeConfig()
const { getCellDefinition } = useCellDefinition()

// Track which zone to highlight
const highlightedZone = ref<string | null>(null)

// Get trident zone distance from env (for approach path highlighting)
const tridentZoneDistance = computed(() => {
    return Number(config.public.tridentZoneDistance) || 1
})

// Mock API response
const scenario = {
    narrative: "A cat jaywalks across the school zone at night, forcing a choice between hitting the animal, colliding with an oncoming bus, or crashing into a high-speed sports car.",
    dilemma_options: {
        maintain: "Maintain: Hit Jaywalking Cat",
        swerve_left: "Swerve Left: Head-on with Bus",
        swerve_right: "Swerve Right: Crash into Sports Car"
    },
    entities: [
        { id: "ent_vehicle_av_ego", type: "vehicle_av", emoji: "🚕", row: 4, col: 14, metadata: { is_star: false, is_ego: true, is_violation: false, action: "", orientation: "W" }},
        { id: "ent_animal_cat_0", type: "animal_cat", emoji: "🐈", row: 4, col: 12, metadata: { is_star: true, is_ego: false, is_violation: true, action: "Jaywalking across road", orientation: "W" }},
        { id: "ent_vehicle_bus_1", type: "vehicle_bus", emoji: "🚌", row: 6, col: 12, metadata: { is_star: false, is_ego: false, is_violation: false, action: "Oncoming at speed", orientation: "E" }},
        { id: "ent_vehicle_sports_car_2", type: "vehicle_sports_car", emoji: "🏎️", row: 3, col: 10, metadata: { is_star: false, is_ego: false, is_violation: false, action: "High-speed approach", orientation: "W" }},
        { id: "ent_vehicle_police_tailgater", type: "vehicle_police", emoji: "🚓", row: 4, col: 15, metadata: { is_star: false, is_ego: false, is_violation: true, action: "Tailgating dangerously close", orientation: "W" }}
    ] as Entity[],
    factors: {
        visibility: "Rain",
        road_condition: "Wet",
        location: "FR",
        brake_status: "Active",
        speed: "High",
        has_tailgater: true,
        primary_entity: "animal_cat",
        primary_behavior: "Violation",
        background_entities: []
    },
    grid_data: [
        [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
        [1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1],
        [3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3],
        [9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9],
        [9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9],
        [12,12,12,12,12,12,12,12,12,12,12,12,12,12,12,12,12,12,12,12],
        [9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9],
        [9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9],
        [4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4],
        [2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2],
        [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]
    ],
    lane_config: {
        E: [[6,0],[6,1],[6,2],[6,3],[6,4],[6,5],[6,6],[6,7],[6,8],[6,9],[6,10],[6,11],[6,12],[6,13],[6,14],[6,15],[6,16],[6,17],[6,18],[6,19],[7,0],[7,1],[7,2],[7,3],[7,4],[7,5],[7,6],[7,7],[7,8],[7,9],[7,10],[7,11],[7,12],[7,13],[7,14],[7,15],[7,16],[7,17],[7,18],[7,19]],
        N: [],
        S: [],
        W: [[3,0],[3,1],[3,2],[3,3],[3,4],[3,5],[3,6],[3,7],[3,8],[3,9],[3,10],[3,11],[3,12],[3,13],[3,14],[3,15],[3,16],[3,17],[3,18],[3,19],[4,0],[4,1],[4,2],[4,3],[4,4],[4,5],[4,6],[4,7],[4,8],[4,9],[4,10],[4,11],[4,12],[4,13],[4,14],[4,15],[4,16],[4,17],[4,18],[4,19]]
    },
    trident_zones: {
        zone_a: {
            coordinates: [
                { row: 4, col: 12, surface: "drivable", orientation: "W" },
                { row: 4, col: 11, surface: "drivable", orientation: "W" },
                { row: 4, col: 10, surface: "drivable", orientation: "W" }
            ]
        },
        zone_b: {
            coordinates: [
                { row: 6, col: 12, surface: "drivable", orientation: "E" },
                { row: 6, col: 11, surface: "drivable", orientation: "E" },
                { row: 6, col: 10, surface: "drivable", orientation: "E" }
            ]
        },
        zone_c: {
            coordinates: [
                { row: 3, col: 12, surface: "drivable", orientation: "W" },
                { row: 3, col: 11, surface: "drivable", orientation: "W" },
                { row: 3, col: 10, surface: "drivable", orientation: "W" }
            ]
        }
    },
    template_name: "Urban School Zone (Straight)",
    current_step: 1,
    total_steps: 12
}

// Get ego vehicle position and orientation
const egoEntity = computed(() => scenario.entities.find(e => e.metadata.is_ego))

// Transform dilemma_options to ranking array with zone mapping
const rankingOptions = computed(() => [
    { key: 'maintain', label: scenario.dilemma_options.maintain, zone: 'zone_a' },
    { key: 'swerve_left', label: scenario.dilemma_options.swerve_left, zone: 'zone_b' },
    { key: 'swerve_right', label: scenario.dilemma_options.swerve_right, zone: 'zone_c' },
])

// Build approach path cells (between ego and zone start)
const approachPathCells = computed(() => {
    if (!egoEntity.value) return new Set<string>()
    
    const ego = egoEntity.value
    const cells = new Set<string>()
    const distance = tridentZoneDistance.value
    
    // Calculate direction deltas based on ego orientation
    let dRow = 0, dCol = 0
    switch (ego.metadata.orientation) {
        case 'N': dRow = -1; break
        case 'S': dRow = 1; break
        case 'E': dCol = 1; break
        case 'W': dCol = -1; break
    }
    
    // Add cells from ego position toward the zone (excluding ego cell itself)
    for (let i = 1; i <= distance; i++) {
        const row = ego.row + (dRow * i)
        const col = ego.col + (dCol * i)
        cells.add(`${row},${col}`)
    }
    
    return cells
})

// Build a lookup map: "row,col" -> highlight type
const zoneCoordMap = computed(() => {
    const map = new Map<string, 'maintain' | 'swerve_left' | 'swerve_right'>()
    
    // zone_a -> maintain
    scenario.trident_zones.zone_a.coordinates.forEach(c => {
        map.set(`${c.row},${c.col}`, 'maintain')
    })
    // zone_b -> swerve_left
    scenario.trident_zones.zone_b.coordinates.forEach(c => {
        map.set(`${c.row},${c.col}`, 'swerve_left')
    })
    // zone_c -> swerve_right
    scenario.trident_zones.zone_c.coordinates.forEach(c => {
        map.set(`${c.row},${c.col}`, 'swerve_right')
    })
    
    return map
})

function getHighlightType(row: number, col: number): 'maintain' | 'swerve_left' | 'swerve_right' | null {
    if (!highlightedZone.value) return null
    
    const cellKey = `${row},${col}`
    const cellType = zoneCoordMap.value.get(cellKey)
    
    // Map zone to highlight type
    const zoneToType: Record<string, 'maintain' | 'swerve_left' | 'swerve_right'> = {
        'zone_a': 'maintain',
        'zone_b': 'swerve_left',
        'zone_c': 'swerve_right'
    }
    
    const activeType = zoneToType[highlightedZone.value]
    if (!activeType) return null
    
    // Highlight zone cells
    if (cellType === activeType) {
        return cellType
    }
    
    // Highlight approach path cells when any zone is highlighted
    if (approachPathCells.value.has(cellKey)) {
        return activeType
    }
    
    return null
}

function getEntitiesAt(row: number, col: number): Entity[] {
    return scenario.entities.filter(e => e.row === row && e.col === col)
}

function getSpeedIcon(s: string): string {
    return { Low: '🐢', Medium: '🚗', High: '🏎️' }[s] || '🚗'
}

function getBrakeIcon(b: string): string {
    return { Active: '✅', Fade: '⚠️', Failed: '🚨' }[b] || '🛞'
}

const speedBadgeClass = computed(() => {
    switch (scenario.factors.speed) {
        case 'Low':
            return 'bg-[hsl(var(--maz-success))]/20 text-[hsl(var(--maz-success))]'
        case 'Medium':
            return 'bg-[hsl(var(--maz-warning))]/20 text-[hsl(var(--maz-warning))]'
        case 'High':
            return 'bg-[hsl(var(--maz-danger))]/20 text-[hsl(var(--maz-danger))]'
        default:
            return 'bg-[hsl(var(--maz-muted))]/20 text-[hsl(var(--maz-muted))]'
    }
})

const brakeBadgeClass = computed(() => {
    switch (scenario.factors.brake_status) {
        case 'Active':
            return 'bg-[hsl(var(--maz-success))]/20 text-[hsl(var(--maz-success))]'
        case 'Fade':
            return 'bg-[hsl(var(--maz-warning))]/20 text-[hsl(var(--maz-warning))]'
        case 'Failed':
            return 'bg-[hsl(var(--maz-danger))]/20 text-[hsl(var(--maz-danger))]'
        default:
            return 'bg-[hsl(var(--maz-muted))]/20 text-[hsl(var(--maz-muted))]'
    }
})
</script>
