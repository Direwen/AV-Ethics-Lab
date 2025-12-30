<template>
    <div class="flex flex-col lg:flex-row gap-6 h-full">
        <!-- Left Column: Board Viewer -->
        <div class="flex-1 flex flex-col items-center">
        <div class="mb-4 w-full max-w-xs">
            <label class="block text-sm font-medium mb-2">Select Board Template</label>
            <select 
            v-model="selectedBoard" 
            class="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 focus:ring-2 focus:ring-primary"
            >
            <option value="board1">Board 1 - School Zone</option>
            <option value="board2">Board 2 - 4-Way Intersection</option>
            <option value="board3">Board 3 - T-Junction</option>
            <option value="custom">Custom (from JSON)</option>
            </select>
        </div>

        <!-- Board Renderer -->
        <div class="flex flex-col items-center p-4 space-y-4">
            <div 
            class="inline-block shadow-2xl border-4 border-gray-900 bg-gray-900 rounded-lg overflow-hidden"
            >
            <div 
                v-for="(row, rIndex) in currentGrid" 
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
            <h2 class="font-bold text-xl">{{ currentName }}</h2>
            <p class="font-mono text-gray-500 text-xs">
                Entities: {{ currentEntities.length }}
            </p>
            <p v-if="narrative" class="text-sm text-gray-600 dark:text-gray-400 max-w-md">
                {{ narrative }}
            </p>
            </div>
        </div>
        </div>

        <!-- Right Column: JSON Input -->
        <div class="flex-1 flex flex-col">
        <label class="block text-sm font-medium mb-2">Scenario JSON Input</label>
        <textarea
            v-model="jsonInput"
            class="flex-1 min-h-[300px] w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 font-mono text-sm focus:ring-2 focus:ring-primary resize-none"
            placeholder="Paste your scenario JSON here..."
        />
        <div class="mt-4 flex gap-2">
            <MazBtn size="sm" @click="loadJson" class="flex-1">
            Load JSON
            </MazBtn>
            <MazBtn size="sm" color="secondary" @click="clearJson" class="flex-1">
            Clear
            </MazBtn>
        </div>
        <p v-if="errorMessage" class="mt-2 text-red-500 text-sm">{{ errorMessage }}</p>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useCellDefinition } from '~/composables/useCellDefinition'
import SimulationCell from '~/components/simulation/Cell.vue'
import type { Entity } from '~/types/simulation'

const { getCellDefinition } = useCellDefinition()

// Entity type to emoji mapping
const ENTITY_EMOJI_MAP: Record<string, string> = {
    'ped_child': '🏃',
    'ped_elderly': '👵',
    'ped_doctor': '👨‍⚕️',
    'ped_criminal': '🦹',
    'vehicle_car': '🚗',
    'vehicle_bus': '🚌',
    'ped_adult': '🧍',
    'obstacle_barrier': '🚧',
}

// Sample board data
const boardData = {
  board1: {
    name: 'Urban School Zone (Straight)',
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
    entities: [] as Entity[]
  },
  board2: {
    name: '4-Way Urban Intersection',
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
    entities: [] as Entity[]
  },
  board3: {
    name: 'Urban T-Junction',
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
    entities: [] as Entity[]
  }
}

const selectedBoard = ref<'board1' | 'board2' | 'board3' | 'custom'>('board1')
const jsonInput = ref('')
const errorMessage = ref('')
const customData = ref<{ grid: number[][], entities: Entity[], name: string, narrative: string } | null>(null)

const currentGrid = computed(() => {
  if (selectedBoard.value === 'custom' && customData.value) {
    return customData.value.grid
  }
  return boardData[selectedBoard.value as keyof typeof boardData]?.grid || []
})

const currentEntities = computed(() => {
  if (selectedBoard.value === 'custom' && customData.value) {
    return customData.value.entities
  }
  return boardData[selectedBoard.value as keyof typeof boardData]?.entities || []
})

const currentName = computed(() => {
  if (selectedBoard.value === 'custom' && customData.value) {
    return customData.value.name
  }
  return boardData[selectedBoard.value as keyof typeof boardData]?.name || ''
})

const narrative = computed(() => {
  if (selectedBoard.value === 'custom' && customData.value) {
    return customData.value.narrative
  }
  return ''
})

function getEntitiesAt(row: number, col: number) {
  return currentEntities.value.filter(
    e => e.position.row === row && e.position.col === col
  )
}

function transformEntity(rawEntity: any): Entity {
  // Use API values directly, fallback to mapping only if missing
  return {
    id: rawEntity.id,
    type: rawEntity.type,
    emoji: rawEntity.emoji || ENTITY_EMOJI_MAP[rawEntity.type] || '❓',
    position: { row: rawEntity.row, col: rawEntity.col },
    metadata: {
      name: rawEntity.metadata?.action || rawEntity.type.replace(/_/g, ' ').replace(/\b\w/g, (c: string) => c.toUpperCase()),
      risk_level: rawEntity.metadata?.is_star ? 'high' : 'low',
      is_occluded: false
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
    
    const entities = (data.entities || []).map(transformEntity)
    
    customData.value = {
      grid: data.grid_data,
      entities,
      name: data.narrative ? 'Custom Scenario' : 'Loaded Scenario',
      narrative: data.narrative || ''
    }
    
    selectedBoard.value = 'custom'
  } catch (e: any) {
    errorMessage.value = e.message || 'Failed to parse JSON'
  }
}

function clearJson() {
  jsonInput.value = ''
  errorMessage.value = ''
  customData.value = null
  selectedBoard.value = 'board1'
}
</script>
