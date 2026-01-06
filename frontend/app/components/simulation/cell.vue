<template>
    <div
        ref="cellRef"
        class="relative w-4 h-3 sm:w-6 sm:h-3 md:w-8 md:h-4 lg:w-12 lg:h-8 flex items-center justify-center transition-all duration-200"
        :class="[
            definition?.class, 
            isHovered && entities?.length && entities.length > 1 ? 'overflow-visible z-50' : '',
            highlightClass,
            roadConditionClass
        ]"
        :title="definition?.name" 
        role="gridcell"
        @mouseenter="isHovered = true"  
        @mouseleave="isHovered = false"
    >
        <!-- Road condition pattern overlay -->
        <div 
            v-if="showRoadConditionOverlay" 
            class="absolute inset-0 pointer-events-none"
            :class="roadConditionOverlayClass"
        />
        <!-- Highlight overlay -->
        <div 
            v-if="highlightType" 
            class="absolute inset-0 pointer-events-none transition-opacity duration-200"
            :class="highlightOverlayClass"
        />
        <EntityRenderer 
            v-if="entities?.length" 
            :entities="entities" 
            :parent-hover="isHovered"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import EntityRenderer from '~/components/simulation/EntityRenderer.vue'
import type { Entity, CellDefinition } from '~/types/simulation'

export type RoadCondition = 'Dry' | 'Wet' | 'Icy'

const props = defineProps<{
    cellCode: number | string
    definition: CellDefinition
    entities?: Entity[]
    highlightType?: 'maintain' | 'swerve_left' | 'swerve_right' | null
    roadCondition?: RoadCondition
}>()

const isHovered = ref(false)

// Check if this cell is a road surface
const isRoadSurface = computed(() => {
    const s = props.definition?.surface
    return s === 'drivable' || s === 'walkable'
})

// Road condition base color shift
const roadConditionClass = computed(() => {
    if (!isRoadSurface.value || !props.roadCondition || props.roadCondition === 'Dry') return ''
    
    switch (props.roadCondition) {
        case 'Wet':
            return 'road-wet-base'
        case 'Icy':
            return 'road-icy-base'
        default:
            return ''
    }
})

// Road condition pattern overlay
const showRoadConditionOverlay = computed(() => {
    console.log(isRoadSurface.value)
    return isRoadSurface.value && props.roadCondition !== 'Dry'
})

const roadConditionOverlayClass = computed(() => {
    if (!props.roadCondition) return ''
    
    switch (props.roadCondition) {
        case 'Wet':
            return 'road-wet-pattern'
        case 'Icy':
            return 'road-icy-pattern'
        default:
            return ''
    }
})

// Highlight classes
const highlightClass = computed(() => {
    if (!props.highlightType) return ''
    return 'ring-2 ring-inset z-10'
})

const highlightOverlayClass = computed(() => {
    switch (props.highlightType) {
        case 'maintain':
            return 'bg-yellow-400/40 ring-yellow-400'
        case 'swerve_left':
            return 'bg-blue-400/40 ring-blue-400'
        case 'swerve_right':
            return 'bg-red-400/40 ring-red-400'
        default:
            return ''
    }
})
</script>
