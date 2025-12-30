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

const { getCellDefinition } = useCellDefinition()

const intersectionData = {
    scenario_id: "uuid-test-002",
    
    template_meta: {
        id: "TPL_002_INTERSECTION",
        name: "4-Way Urban Intersection",
        dimensions: { rows: 11, cols: 20 }
    },

    grid: [
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 14, 10, 20, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 14, 10, 20, 0, 0, 0, 0, 0, 0, 0],
        [1, 1, 1, 1, 1, 1, 1, 1, 19, 15, 15, 15, 20, 1, 1, 1, 1, 1, 1, 1],
        [3, 3, 3, 3, 3, 3, 3, 3, 5, 11, 11, 11, 7, 3, 3, 3, 3, 3, 3, 3],
        [9, 9, 9, 9, 9, 9, 9, 16, 11, 11, 11, 11, 11, 16, 9, 9, 9, 9, 9, 9],
        [13, 13, 13, 13, 13, 13, 13, 16, 11, 11, 11, 11, 11, 16, 13, 13, 13, 13, 13, 13],
        [9, 9, 9, 9, 9, 9, 9, 16, 11, 11, 11, 11, 11, 16, 9, 9, 9, 9, 9, 9],
        [4, 4, 4, 4, 4, 4, 4, 4, 6, 11, 11, 11, 8, 4, 4, 4, 4, 4, 4, 4],
        [2, 2, 2, 2, 2, 2, 2, 2, 19, 15, 15, 15, 20, 2, 2, 2, 2, 2, 2, 2],
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 14, 10, 20, 0, 0, 0, 0, 0, 0, 0],
        [0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 14, 10, 20, 0, 0, 0, 0, 0, 0, 0]
    ],

    entities: [
        {
            id: "ent_mother_01",
            type: "pedestrian",
            emoji: "🤱",
            position: { row: 3, col: 1 },
            metadata: {
                name: "Mother with Baby",
                risk_level: "high",
                is_occluded: false
            }
        },
        {
            id: "ent_car_01",
            type: "vehicle",
            emoji: "🚗",
            position: { row: 5, col: 3 },
            metadata: {
                name: "Sedan",
                risk_level: "none",
                is_occluded: false
            }
        },
        {
            id: "ent_bike_01",
            type: "vehicle",
            emoji: "🚴",
            position: { row: 4, col: 15 },
            metadata: {
                name: "Cyclist",
                risk_level: "medium",
                is_occluded: false
            }
        },
        {
            id: "ent_ped_01",
            type: "pedestrian",
            emoji: "🚶",
            position: { row: 7, col: 13 },
            metadata: {
                name: "Pedestrian Crossing",
                risk_level: "high",
                is_occluded: false
            }
        }
    ]
}

const currentTemplate = computed(() => intersectionData)

function getEntitiesAt(row: number, col: number) {
    return currentTemplate.value.entities.filter(
        e => e.position.row === row && e.position.col === col
    )
}
</script>
