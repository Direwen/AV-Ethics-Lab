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

// --- 1. MOCK DATA (Backend Structure) ---
const schoolZoneData = {
  scenario_id: "uuid-test-001",
  
  template_meta: {
    id: "TPL_001_SCHOOL_ZONE",
    name: "Urban School Zone (Straight)",
    dimensions: { rows: 11, cols: 20 }
  },
  
  // The Map (Visuals)
  grid: [
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3],
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9], // Top Lane
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],
    [12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12], // Dash Line
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9],
    [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9], // Bottom Lane
    [4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4],
    [2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
  ],

  // The Entities (Dynamic)
  entities: [
    // Case 1: Single Child on Sidewalk (Top)
    {
        id: "ent_01",
        type: "child",
        emoji: "🧍",
        position: { row: 2, col: 5 },
        metadata: {
            name: "Child (Waiting)",
            risk_level: "low",
            is_occluded: false
        }
    },
    // Case 2: The "Stack" (Bus + Hidden Child) in the Road
    {
        id: "ent_bus_01",
        type: "vehicle",
        emoji: "🚌",
        position: { row: 4, col: 10 },
        metadata: {
            name: "School Bus",
            risk_level: "none",
            is_occluded: false
        }
    },
    {
        id: "ent_child_hidden",
        type: "child",
        emoji: "🧍",
        position: { row: 4, col: 10 }, // SAME COORDINATES
        metadata: {
            name: "Child (Occluded)",
            risk_level: "critical",
            is_occluded: true // Should appear greyed out or behind
        }
    },
    {
        id: "ent_child_hidden1",
        type: "child",
        emoji: "🧍",
        position: { row: 4, col: 10 }, // SAME COORDINATES
        metadata: {
            name: "Child (Occluded)",
            risk_level: "critical",
            is_occluded: true // Should appear greyed out or behind
        }
    },
    {
        id: "ent_child_hidden2",
        type: "child",
        emoji: "🧍",
        position: { row: 4, col: 10 }, // SAME COORDINATES
        metadata: {
            name: "Child (Occluded)",
            risk_level: "critical",
            is_occluded: true // Should appear greyed out or behind
        }
    },
    // Case 3: Ego Vehicle (You)
    {
        id: "ent_ego",
        type: "vehicle",
        emoji: "🚗",
        position: { row: 4, col: 2 },
        metadata: {
            name: "Ego Vehicle",
            risk_level: "none",
            is_occluded: false
        }
    }
  ]
}

const currentTemplate = computed(() => schoolZoneData)

// --- 2. HELPER: FILTER ENTITIES ---
function getEntitiesAt(row: number, col: number) {
  return currentTemplate.value.entities.filter(
    e => e.position.row === row && e.position.col === col
  )
}
</script>