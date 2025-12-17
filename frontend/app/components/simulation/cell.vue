<template>
    <div
        class="hover:bg-green-200/50 relative w-4 h-4 sm:w-6 sm:h-6 md:w-8 md:h-8 lg:w-12 lg:h-12 flex items-center justify-center transition-all duration-30"
        :class="[
            definition?.class, // The Visuals (e.g., 'bg-gray-800')
            isInteractive ? 'hover:brightness-110 cursor-pointer' : '' // Interaction hints
        ]"
        :title="definition?.name" 
        role="gridcell"
        :aria-label="definition?.name"
        @click="handleClick"
    >
        <EntityRenderer v-if="entities?.length" :entities="entities" />
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import EntityRenderer from '~/components/simulation/EntityRenderer.vue'
import type { Entity } from '~/types/simulation'

// Define the shape of your Metadata based on your JSON structure
export interface CellDefinition {
    name?: string
    class: string
    allow: string[]
    risk_factor?: string
}

const props = defineProps<{
    cellCode: number | string
    definition: CellDefinition
    entities?: Entity[]
    // Future-proofing: Is this cell valid for the currently selected item?
    isValidTarget?: boolean 
}>()

const emit = defineEmits(['cell-click'])

// Computed Logic
const isInteractive = computed(() => {
    // A cell is interactive if it allows *something* (i.e., not a building roof)
    return props.definition?.allow?.length > 0
})

function handleClick() {
    if (!isInteractive.value) return
    
    // Emit the full context back to the Board
    emit('cell-click', {
        code: props.cellCode,
        allowed: props.definition.allow,
        type: props.definition.name
    })
}
</script>