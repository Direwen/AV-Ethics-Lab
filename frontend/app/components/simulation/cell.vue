<template>
    <div
        ref="cellRef"
        class="relative w-4 h-3 sm:w-6 sm:h-3 md:w-8 md:h-4 lg:w-12 lg:h-8 flex items-center justify-center transition-all duration-30"
        :class="[
            definition?.class, 
            isHovered && entities?.length && entities.length > 1 ? 'overflow-visible z-50' : ''
        ]"
        :title="definition?.name" 
        role="gridcell"
        @mouseenter="isHovered = true"  
        @mouseleave="isHovered = false"
    >
        <EntityRenderer 
            v-if="entities?.length" 
            :entities="entities" 
            :parent-hover="isHovered"
        />
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import EntityRenderer from '~/components/simulation/EntityRenderer.vue'
import type { Entity, CellDefinition } from '~/types/simulation'

defineProps<{
    cellCode: number | string
    definition: CellDefinition
    entities?: Entity[]
}>()

const isHovered = ref(false)
</script>
