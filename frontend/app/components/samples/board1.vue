<template>
  <div class="flex flex-col items-center p-4 sm:p-6 md:p-10 space-y-6">
    <div class="inline-block shadow-2xl border-4 border-gray-900 bg-gray-900 rounded-lg overflow-hidden">
      <div v-for="(row, rIndex) in templateData.grid" :key="rIndex" class="flex">
        <div v-for="(cellCode, cIndex) in row" :key="cIndex" class="relative">
          <SimulationCell
            :cell-code="cellCode"
            :definition="getCellDefinition(String(cellCode))"
            :entities="getEntitiesAt(rIndex, cIndex)"
          />
          <div 
            v-if="getLaneDirection(rIndex, cIndex)"
            class="absolute inset-0 flex items-center justify-center pointer-events-none"
          >
            <span class="text-xs font-bold" :class="getLaneArrowClass(getLaneDirection(rIndex, cIndex))">
              {{ getLaneArrow(getLaneDirection(rIndex, cIndex)) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <LaneDirectionLegend :directions="activeDirections" />

    <div class="text-center space-y-2">
      <h2 class="font-bold text-xl">{{ templateData.name }}</h2>
      <p class="font-mono text-gray-500 text-xs">
        Template ID: {{ templateData.meta.id }} <br/>
        Dimensions: {{ templateData.width }} x {{ templateData.height }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useCellDefinition } from '~/composables/useCellDefinition'
import { useLaneDirection } from '~/composables/useLaneDirection'

const { getCellDefinition } = useCellDefinition()

const templateData = {
  name: "Urban School Zone (Straight)",
  width: 20,
  height: 11,
  grid: [
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3],
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],
    [12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12],
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],
    [4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4],
    [2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
  ],
  laneConfig: {
    W: [[3,0],[3,1],[3,2],[3,3],[3,4],[3,5],[3,6],[3,7],[3,8],[3,9],[3,10],[3,11],[3,12],[3,13],[3,14],[3,15],[3,16],[3,17],[3,18],[3,19],
        [4,0],[4,1],[4,2],[4,3],[4,4],[4,5],[4,6],[4,7],[4,8],[4,9],[4,10],[4,11],[4,12],[4,13],[4,14],[4,15],[4,16],[4,17],[4,18],[4,19]],
    E: [[6,0],[6,1],[6,2],[6,3],[6,4],[6,5],[6,6],[6,7],[6,8],[6,9],[6,10],[6,11],[6,12],[6,13],[6,14],[6,15],[6,16],[6,17],[6,18],[6,19],
        [7,0],[7,1],[7,2],[7,3],[7,4],[7,5],[7,6],[7,7],[7,8],[7,9],[7,10],[7,11],[7,12],[7,13],[7,14],[7,15],[7,16],[7,17],[7,18],[7,19]],
    N: [],
    S: []
  },
  meta: { id: "TPL_001_SCHOOL_ZONE" },
  entities: [
    { id: "ent_ego", type: "vehicle", emoji: "🚗", position: { row: 4, col: 2 }, metadata: { name: "Ego Vehicle (Westbound)", risk_level: "none", is_occluded: false } },
    { id: "ent_bus_01", type: "vehicle", emoji: "🚌", position: { row: 7, col: 15 }, metadata: { name: "School Bus (Eastbound)", risk_level: "none", is_occluded: false } }
  ]
}

const laneConfig = computed(() => templateData.laneConfig)
const { getLaneDirection, getLaneArrow, getLaneArrowClass, activeDirections } = useLaneDirection(laneConfig)

function getEntitiesAt(row: number, col: number) {
  return templateData.entities.filter(e => e.position.row === row && e.position.col === col)
}
</script>
