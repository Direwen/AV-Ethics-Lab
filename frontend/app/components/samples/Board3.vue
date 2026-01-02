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
  name: "Urban T-Junction",
  width: 20,
  height: 11,
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
  laneConfig: {
    W: [[7,0],[7,1],[7,2],[7,3],[7,4],[7,5],[7,6],[7,7],[7,8],[7,9],[7,10],[7,11],[7,12],[7,13],[7,14],[7,15],[7,16],[7,17],[7,18],[7,19]],
    E: [[9,0],[9,1],[9,2],[9,3],[9,4],[9,5],[9,6],[9,7],[9,8],[9,9],[9,10],[9,11],[9,12],[9,13],[9,14],[9,15],[9,16],[9,17],[9,18],[9,19]],
    S: [[0,9],[1,9],[2,9],[3,9],[4,9],[5,9],[6,9],[7,9],[8,9],[9,9]],
    N: [[0,11],[1,11],[2,11],[3,11],[4,11],[5,11],[6,11],[7,11],[8,11],[9,11]]
  },
  meta: { id: "TPL_003_T_JUNCTION" },
  entities: [
    { id: "ent_car_01", type: "vehicle", emoji: "🚗", position: { row: 2, col: 9 }, metadata: { name: "Car (Southbound)", risk_level: "none", is_occluded: false } },
    { id: "ent_truck_01", type: "vehicle", emoji: "🚚", position: { row: 7, col: 4 }, metadata: { name: "Truck (Westbound)", risk_level: "none", is_occluded: false } },
    { id: "ent_car_02", type: "vehicle", emoji: "🚙", position: { row: 9, col: 16 }, metadata: { name: "SUV (Eastbound)", risk_level: "none", is_occluded: false } },
    { id: "ent_bus_01", type: "vehicle", emoji: "🚌", position: { row: 6, col: 11 }, metadata: { name: "Bus (Northbound)", risk_level: "none", is_occluded: false } }
  ]
}

const laneConfig = computed(() => templateData.laneConfig)
const { getLaneDirection, getLaneArrow, getLaneArrowClass, activeDirections } = useLaneDirection(laneConfig)

function getEntitiesAt(row: number, col: number) {
  return templateData.entities.filter(e => e.position.row === row && e.position.col === col)
}
</script>
