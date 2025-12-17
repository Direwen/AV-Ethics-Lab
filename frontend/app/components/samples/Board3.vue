<template>
    <div class="flex flex-col items-center p-4 sm:p-6 md:p-10 space-y-6">
        
        <div 
        class="inline-block shadow-2xl border-4 border-gray-900 bg-gray-900 rounded-lg overflow-hidden"
        >
            <div 
                v-for="(row, rIndex) in currentTemplate.grid" 
                :key="rIndex" 
                class="flex"
            >
                <SimulationCell
                v-for="(cellCode, cIndex) in row" 
                :key="cIndex"
                :cell-code="cellCode"
                :definition="getCellDefinition(String(cellCode))"
                >
                    <EntityRenderer 
                        v-if="getEntitiesAt(rIndex, cIndex).length > 0"
                        :entities="getEntitiesAt(rIndex, cIndex)" 
                    />
                </SimulationCell>
            </div>
        </div>

        <div class="text-center space-y-2">
            <h2 class="font-bold text-xl">{{ currentTemplate.template_meta.name }}</h2>
            <p class="font-mono text-gray-500 text-xs">
                Scenario ID: {{ currentTemplate.scenario_id }} <br/>
                Resolution: {{ currentTemplate.template_meta.dimensions.cols }}x{{ currentTemplate.template_meta.dimensions.rows }}
            </p>
        </div>

    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useCellDefinition } from '~/composables/useCellDefinition' // Or your constant import
import SimulationCell from '~/components/simulation/Cell.vue'
import EntityRenderer from '~/components/simulation/EntityRenderer.vue'

// Use the composable or import CELL_DEFINITIONS directly depending on your setup
const { getCellDefinition } = useCellDefinition()

// --- MOCK DATA FOR T-JUNCTION ---
const tJunctionData = {
    scenario_id: "uuid-test-003",
    
    template_meta: {
        id: "TPL_003_T_JUNCTION",
        name: "Urban T-Junction",
        dimensions: { rows: 11, cols: 20 }
    },

    // GRID CODES:
    // 0=Roof, 2=Building Edge Bottom
    // 10=Vert Road, 9=Horiz Road
    // 4=Sidewalk Bottom, 3=Sidewalk Top
    // 6=Corner BR, 8=Corner BL
    // 16=Crosswalk/Stop Line, 12=Dash Line
    grid: [
        // Rows 0-4: Buildings (Left/Right) & Vertical Road (Center)
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 17, 10, 20, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 17, 10, 20, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 17, 10, 20, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 17, 10, 20, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 15, 15, 15, 20, 0, 0, 0, 0, 0, 0, 0],
        
        // Row 5: Building Edge (Visual Depth)
        [2, 2, 2, 2, 2, 2, 2, 8, 19, 10, 17, 10, 20, 2, 2, 2, 2, 2, 2, 2],

        // Row 6: Sidewalks & Corners turning into Junction
        // Code 6 (Corner BR) and 8 (Corner BL) create the curve
        [4, 4, 4, 4, 4, 4, 4, 4, 5, 10, 17, 10, 7, 4, 4, 4, 4, 4, 4, 4],

        // Row 7: Horizontal Road Top Lane
        [9, 9, 9, 9, 9, 9, 9, 9, 11, 11, 11, 11, 9, 9, 9, 9, 9, 9, 9, 9],

        // Row 8: Center Yellow Line
        [18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18],

        // Row 9: Horizontal Road Bottom Lane
        [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],

        // Row 10: Bottom Sidewalk
        [3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3]
    ],

  // ENTITIES (Example: Car at Stop Line, Pedestrian on Corner)
    entities: [
        {
            id: "ent_car_01",
            type: "vehicle",
            emoji: "🚗",
            position: { row: 5, col: 9 }, // Driving down vertical road
            metadata: {
                name: "Approaching Vehicle",
                risk_level: "medium",
                is_occluded: false
            }
        },
        {
            id: "ent_ped_01",
            type: "child",
            emoji: "🏃",
            position: { row: 6, col: 7 }, // On the corner sidewalk
            metadata: {
                name: "Child (Near Road)",
                risk_level: "high",
                is_occluded: false
            }
        }
    ]
}

const currentTemplate = computed(() => tJunctionData)

// Helper: Filter entities
function getEntitiesAt(row: number, col: number) {
    return currentTemplate.value.entities.filter(
        e => e.position.row === row && e.position.col === col
    )
}
</script>