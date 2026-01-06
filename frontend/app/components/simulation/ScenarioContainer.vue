<template>
    <div class="relative inline-block">
        <!-- Main content with visibility filter -->
        <div 
            class="shadow-2xl border-4 border-gray-900 bg-gray-900 rounded-lg overflow-hidden"
            :class="visibilityClass"
        >
            <slot />
        </div>
        
        <!-- Visibility overlay (fog haze, rain effect) -->
        <div 
            v-if="showOverlay" 
            class="absolute inset-0 pointer-events-none rounded-lg"
            :class="overlayClass"
        />
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export type Visibility = 'Clear' | 'Fog' | 'Night' | 'Rain'

const props = defineProps<{
    visibility?: Visibility
}>()

const visibilityClass = computed(() => {
    switch (props.visibility) {
        case 'Night':
            return 'visibility-night'
        case 'Fog':
            return 'visibility-fog'
        case 'Rain':
            return 'visibility-rain'
        default:
            return ''
    }
})

const showOverlay = computed(() => {
    return props.visibility === 'Fog' || props.visibility === 'Rain'
})

const overlayClass = computed(() => {
    switch (props.visibility) {
        case 'Fog':
            return 'visibility-fog-overlay'
        case 'Rain':
            return 'visibility-rain-overlay'
        default:
            return ''
    }
})
</script>
