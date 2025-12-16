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
          v-for="(cell, idx) in row" 
          :key="idx"
          :cell-code="cell"
          :definition="getCellDefinition(cell)"
        />
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

const schoolZoneTemplate: {
  id: string
  name: string
  dimensions: { rows: number; cols: number }
  grid: number[][]
} = {
  id: "TPL_001_SCHOOL_ZONE",
  name: "Urban School Zone (Straight)",
  dimensions: { rows: 11, cols: 20 },
  grid: [
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3],
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],
    [9, 9, 9, 12, 9, 9, 9, 12, 9, 9, 9, 12, 9, 9, 9, 12, 9, 9, 9, 12],
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],
    [4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4],
    [2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
  ]
}

const currentTemplate = computed(() => schoolZoneTemplate)
</script>