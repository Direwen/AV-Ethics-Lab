<template>
    <div class="flex flex-col items-center p-4 sm:p-6 md:p-10 space-y-6">
        <div class="inline-block shadow-2xl border-4 border-gray-900 bg-gray-900 rounded-lg overflow-hidden">
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
                    :entities="getEntitiesAt(rIndex, cIndex)"
                />
            </div>
        </div>

        <div class="text-center space-y-2">
            <h2 class="font-bold text-xl">{{ currentTemplate.template_meta.name }}</h2>
            <p class="font-mono text-gray-500 text-xs">
                Scenario ID: {{ currentTemplate.scenario_id }} <br/>
                Entities: {{ currentTemplate.entities.length }}
            </p>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useCellDefinition } from '~/composables/useCellDefinition'
import SimulationCell from '~/components/simulation/Cell.vue'

const { getCellDefinition } = useCellDefinition()

const tJunctionData = {
    scenario_id: "uuid-test-003",
    
    template_meta: {
        id: "TPL_003_T_JUNCTION",
        name: "Urban T-Junction",
        dimensions: { rows: 11, cols: 20 }
    },

    grid: [
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 17, 10, 20, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 17, 10, 20, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 17, 10, 20, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 17, 10, 20, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 15, 15, 15, 20, 0, 0, 0, 0, 0, 0, 0],
        [2, 2, 2, 2, 2, 2, 2, 8, 19, 10, 17, 10, 20, 2, 2, 2, 2, 2, 2, 2],
        [4, 4, 4, 4, 4, 4, 4, 4, 5, 10, 17, 10, 7, 4, 4, 4, 4, 4, 4, 4],
        [9, 9, 9, 9, 9, 9, 9, 9, 11, 11, 11, 11, 9, 9, 9, 9, 9, 9, 9, 9],
        [18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18, 18],
        [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],
        [3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3]
    ],

    entities: [
        {
            id: "ent_car_01",
            type: "vehicle",
            emoji: "🚗",
            position: { row: 3, col: 9 },
            metadata: {
                name: "Approaching Vehicle",
                risk_level: "medium",
                is_occluded: false
            }
        },
        {
            id: "ent_child_01",
            type: "child",
            emoji: "🏃",
            position: { row: 6, col: 7 },
            metadata: {
                name: "Child (Near Road)",
                risk_level: "high",
                is_occluded: false
            }
        },
        {
            id: "ent_truck_01",
            type: "vehicle",
            emoji: "🚚",
            position: { row: 7, col: 4 },
            metadata: {
                name: "Delivery Truck",
                risk_level: "none",
                is_occluded: false
            }
        },
        {
            id: "ent_dog_01",
            type: "animal",
            emoji: "🐕",
            position: { row: 10, col: 12 },
            metadata: {
                name: "Stray Dog",
                risk_level: "low",
                is_occluded: false
            }
        }
    ]
}

const currentTemplate = computed(() => tJunctionData)

function getEntitiesAt(row: number, col: number) {
    return currentTemplate.value.entities.filter(
        e => e.position.row === row && e.position.col === col
    )
}
</script>
