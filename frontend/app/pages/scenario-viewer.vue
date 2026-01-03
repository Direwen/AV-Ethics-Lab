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
        <!-- Factors Badge Panel -->
        <div v-if="customData?.factors" class="flex flex-wrap gap-2 justify-center">
          <span class="px-3 py-1.5 rounded-full text-xs font-medium bg-indigo-100 dark:bg-indigo-900/50 text-indigo-800 dark:text-indigo-200 flex items-center gap-1">
            {{ getVisibilityIcon(customData.factors.visibility) }} {{ customData.factors.visibility }}
          </span>
          <span class="px-3 py-1.5 rounded-full text-xs font-medium bg-amber-100 dark:bg-amber-900/50 text-amber-800 dark:text-amber-200 flex items-center gap-1">
            {{ getRoadConditionIcon(customData.factors.road_condition) }} {{ customData.factors.road_condition }}
          </span>
          <span class="px-3 py-1.5 rounded-full text-xs font-medium bg-green-100 dark:bg-green-900/50 text-green-800 dark:text-green-200 flex items-center gap-1">
            {{ getSpeedIcon(customData.factors.speed) }} {{ customData.factors.speed }}
          </span>
          <span 
            class="px-3 py-1.5 rounded-full text-xs font-medium flex items-center gap-1"
            :class="customData.factors.brake_status === 'Active' 
              ? 'bg-green-100 dark:bg-green-900/50 text-green-800 dark:text-green-200'
              : 'bg-red-100 dark:bg-red-900/50 text-red-800 dark:text-red-200'"
          >
            {{ getBrakeIcon(customData.factors.brake_status) }} Brake: {{ customData.factors.brake_status }}
          </span>
          <span class="px-3 py-1.5 rounded-full text-xs font-medium bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 flex items-center gap-1">
            📍 {{ customData.factors.location }}
          </span>
          <span 
            v-if="customData.factors.has_tailgater"
            class="px-3 py-1.5 rounded-full text-xs font-medium bg-orange-100 dark:bg-orange-900/50 text-orange-800 dark:text-orange-200 flex items-center gap-1"
          >
            🚙 Tailgater
          </span>
        </div>

        <!-- Map with visual overlay -->
        <div 
          v-if="customData"
          class="inline-block shadow-2xl border-4 border-gray-900 bg-gray-900 rounded-lg overflow-hidden transition-all duration-300"
          :class="mapOverlayClasses"
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
    
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useCellDefinition } from '~/composables/useCellDefinition'
import { useLaneDirection, type LaneConfig } from '~/composables/useLaneDirection'
import type { Entity } from '~/types/simulation'

const { getCellDefinition } = useCellDefinition()

interface ScenarioFactors {
  visibility: string
  road_condition: string
  location: string
  brake_status: string
  speed: string
  has_tailgater: boolean
  primary_entity: string
  primary_behavior: string
  background_entities: string[]
}

interface CustomScenarioData {
  grid: number[][]
  entities: Entity[]
  narrative: string
  width: number
  height: number
  laneConfig: LaneConfig
  factors: ScenarioFactors | null
}

const jsonInput = ref('')
const errorMessage = ref('')
const customData = ref<CustomScenarioData | null>(null)

const laneConfig = computed(() => customData.value?.laneConfig || { W: [], E: [], N: [], S: [] })
const { getLaneDirection, getLaneArrow, getLaneArrowClass, activeDirections } = useLaneDirection(laneConfig)

// Visual overlay classes based on factors
const mapOverlayClasses = computed(() => {
  const factors = customData.value?.factors
  if (!factors) return ''
  
  const classes: string[] = []
  
  // Visibility effects (matching backend: Clear, Fog, Night, Rain)
  if (factors.visibility === 'Night') classes.push('brightness-[0.6]')
  else if (factors.visibility === 'Fog') classes.push('brightness-90 blur-[0.5px]')
  else if (factors.visibility === 'Rain') classes.push('brightness-[0.85]')
  // Clear = no effect
  
  // Road condition effects (matching backend: Dry, Wet, Icy)
  if (factors.road_condition === 'Wet') classes.push('contrast-[0.95] saturate-110')
  else if (factors.road_condition === 'Icy') classes.push('saturate-50 brightness-105')
  // Dry = no effect
  
  return classes.join(' ')
})

// Icon helpers for factors (matching backend constants)
function getVisibilityIcon(visibility: string): string {
  const icons: Record<string, string> = {
    'Clear': '☀️',
    'Fog': '�️',
    'Night': '🌙',
    'Rain': '🌧️'
  }
  return icons[visibility] || '👁️'
}

function getRoadConditionIcon(condition: string): string {
  const icons: Record<string, string> = {
    'Dry': '🛣️',
    'Wet': '💧',
    'Icy': '🧊'
  }
  return icons[condition] || '🛣️'
}

function getSpeedIcon(speed: string): string {
  const icons: Record<string, string> = {
    'Low': '🐢',
    'Medium': '🚕',
    'High': '🏎️'
  }
  return icons[speed] || '🚕'
}

function getBrakeIcon(status: string): string {
  const icons: Record<string, string> = {
    'Active': '✅',
    'Fade': '⚠️',
    'Failed': '🚨'
  }
  return icons[status] || '🛞'
}

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
      laneConfig: data.lane_config || { W: [], E: [], N: [], S: [] },
      factors: data.factors || null
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
