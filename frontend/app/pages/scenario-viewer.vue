<template>
  <div class="flex flex-col gap-6 h-full">

    <div class="flex-1 flex flex-col">
      <label class="block text-sm font-medium mb-2">Scenario JSON Input</label>
      <textarea
        v-model="jsonInput"
        class="flex-1 min-h-[300px] w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 font-mono text-sm focus:ring-2 focus:ring-primary resize-none"
        placeholder="Paste your scenario JSON here..."
      />
      <div class="mt-4 flex gap-2">
        <MazBtn size="sm" @click="loadJson" class="flex-1">Load JSON</MazBtn>
        <MazBtn size="sm" color="secondary" @click="clearJson" class="flex-1">Clear</MazBtn>
      </div>
      <p v-if="errorMessage" class="mt-2 text-red-500 text-sm">{{ errorMessage }}</p>
    </div>

    <div class="flex-1 flex flex-col items-center">
      <div class="flex flex-col items-center p-4 space-y-4">
        <div 
          v-if="customData"
          class="inline-block shadow-2xl border-4 border-gray-900 bg-gray-900 rounded-lg overflow-hidden"
        >
          <div v-for="(row, rIndex) in customData.grid" :key="rIndex" class="flex">
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

        <!-- Empty State -->
        <div 
          v-else 
          class="flex flex-col items-center justify-center p-12 border-2 border-dashed border-gray-400 dark:border-gray-600 rounded-lg"
        >
          <span class="text-4xl mb-4">🗺️</span>
          <p class="text-gray-500 dark:text-gray-400 text-center">
            Paste scenario JSON and click "Load JSON" to render the map
          </p>
        </div>

        <LaneDirectionLegend v-if="customData" :directions="activeDirections" />

        <div v-if="customData" class="text-center space-y-2">
          <h2 class="font-bold text-xl">Custom Scenario</h2>
          <p class="font-mono text-gray-500 text-xs">
            {{ customData.width }} x {{ customData.height }} | Entities: {{ customData.entities.length }}
          </p>
          <p v-if="customData.narrative" class="text-sm text-gray-600 dark:text-gray-400 max-w-md">
            {{ customData.narrative }}
          </p>
        </div>
      </div>
    </div>

    <!-- Right Column: JSON Input -->
    
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useCellDefinition } from '~/composables/useCellDefinition'
import { useLaneDirection, type LaneConfig } from '~/composables/useLaneDirection'
import type { Entity } from '~/types/simulation'

const { getCellDefinition } = useCellDefinition()

interface CustomScenarioData {
  grid: number[][]
  entities: Entity[]
  narrative: string
  width: number
  height: number
  laneConfig: LaneConfig
}

const jsonInput = ref('')
const errorMessage = ref('')
const customData = ref<CustomScenarioData | null>(null)

const laneConfig = computed(() => customData.value?.laneConfig || { W: [], E: [], N: [], S: [] })
const { getLaneDirection, getLaneArrow, getLaneArrowClass, activeDirections } = useLaneDirection(laneConfig)

function getEntitiesAt(row: number, col: number) {
  if (!customData.value) return []
  return customData.value.entities.filter(e => e.row === row && e.col === col)
}

function transformEntity(rawEntity: any): Entity {
  return {
    id: rawEntity.id,
    type: rawEntity.type,
    emoji: rawEntity.emoji || '❓',
    row: rawEntity.row,
    col: rawEntity.col,
    metadata: {
      is_star: rawEntity.metadata?.is_star ?? false,
      is_ego: rawEntity.metadata?.is_ego ?? false,
      is_violation: rawEntity.metadata?.is_violation ?? false,
      action: rawEntity.metadata?.action ?? '',
      orientation: rawEntity.metadata?.orientation ?? ''
    }
  }
}

function loadJson() {
  errorMessage.value = ''
  try {
    const parsed = JSON.parse(jsonInput.value)
    const data = parsed.data || parsed
    
    if (!data.grid_data || !Array.isArray(data.grid_data)) {
      throw new Error('Invalid JSON: missing grid_data array')
    }
    
    customData.value = {
      grid: data.grid_data,
      entities: (data.entities || []).map(transformEntity),
      narrative: data.narrative || '',
      width: data.width || data.grid_data[0]?.length || 0,
      height: data.height || data.grid_data.length,
      laneConfig: data.lane_config || { W: [], E: [], N: [], S: [] }
    }
  } catch (e: any) {
    errorMessage.value = e.message || 'Failed to parse JSON'
  }
}

function clearJson() {
  jsonInput.value = ''
  errorMessage.value = ''
  customData.value = null
}
</script>
