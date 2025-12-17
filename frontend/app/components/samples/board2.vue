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
                    :definition="getCellDefinition(cellCode)"
                    @cell-click=""
                >
                    <span 
                        v-if="rIndex === 3 && cIndex === 1" 
                        class="text-xs sm:text-sm md:text-2lg lg:text-5xl"
                    >🤱</span>
                </SimulationCell>
            </div>
        </div>

        <div class="text-center space-y-2">
        <h2 class="font-bold text-xl">{{ currentTemplate.name }}</h2>
        <p class="font-mono text-gray-500 text-xs">
            Resolution: {{ currentTemplate.dimensions.cols }}x{{ currentTemplate.dimensions.rows }} | 
            ID: {{ currentTemplate.id }}
        </p>
        </div>

    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useCellDefinition } from '~/composables/useCellDefinition'

const { getCellDefinition } = useCellDefinition()

const intersectionTemplate = {
    id: "TPL_002_INTERSECTION",
    name: "4-Way Urban Intersection",
    dimensions: { rows: 11, cols: 20 },
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
    ]
}

const currentTemplate = computed(() => intersectionTemplate)
</script>