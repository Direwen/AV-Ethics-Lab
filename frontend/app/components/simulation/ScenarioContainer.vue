<template>
    <div ref="wrapperRef" class="w-full flex justify-center items-start"> <div 
            class="relative inline-block origin-top transition-transform duration-150"
            :style="{ transform: `scale(${scale})` }"
        >
            <div 
                ref="contentRef"
                class="shadow-2xl border-4 border-gray-900 bg-gray-900 rounded-lg overflow-hidden block" 
                :class="visibilityClass"
            >
                <slot />
            </div>
            
            <div 
                v-if="showOverlay" 
                class="absolute inset-0 pointer-events-none rounded-lg"
                :class="overlayClass"
            />
        </div>
        <div :style="{ height: `${scaledHeight}px`, marginTop: `-${scaledHeight}px` }" class="pointer-events-none w-px" />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'

export type Visibility = 'Clear' | 'Fog' | 'Night' | 'Rain'

const props = withDefaults(defineProps<{
    visibility?: Visibility
    // New Prop: Allow scaling up to 1.5x (or more) of original size
    maxScale?: number 
}>(), {
    maxScale: 1.5 // Default to allowing 50% growth
})

const wrapperRef = ref<HTMLElement | null>(null)
const contentRef = ref<HTMLElement | null>(null)
const containerWidth = ref(0)
const contentWidth = ref(0)
const contentHeight = ref(0)

// Scale factor: 
// 1. Calculate ratio (Available Space / Actual Content)
// 2. Cap it at maxScale (so it doesn't get pixelated on huge screens)
const scale = computed(() => {
    if (!contentWidth.value || !containerWidth.value) return 1
    const ratio = containerWidth.value / contentWidth.value
    // The Magic Fix: We allow it to go above 1.0, up to maxScale
    return Math.min(props.maxScale, ratio)
})

// Adjusted height after scaling
const scaledHeight = computed(() => {
    return contentHeight.value * scale.value
})

let resizeObserver: ResizeObserver | null = null

function updateDimensions() {
    if (wrapperRef.value) {
        containerWidth.value = wrapperRef.value.clientWidth
    }
    if (contentRef.value) {
        // Use offsetWidth/Height for more accurate box-model measurement
        contentWidth.value = contentRef.value.offsetWidth
        contentHeight.value = contentRef.value.offsetHeight
    }
}

onMounted(() => {
    nextTick(() => {
        updateDimensions()
        
        resizeObserver = new ResizeObserver(() => {
            updateDimensions()
        })
        
        if (wrapperRef.value) {
            resizeObserver.observe(wrapperRef.value)
        }
        if (contentRef.value) {
            resizeObserver.observe(contentRef.value)
        }
    })
})

onUnmounted(() => {
    resizeObserver?.disconnect()
})

const visibilityClass = computed(() => {
    switch (props.visibility) {
        case 'Night': return 'visibility-night'
        case 'Fog': return 'visibility-fog'
        case 'Rain': return 'visibility-rain'
        default: return ''
    }
})

const showOverlay = computed(() => {
    return props.visibility === 'Fog' || props.visibility === 'Rain'
})

const overlayClass = computed(() => {
    switch (props.visibility) {
        case 'Fog': return 'visibility-fog-overlay'
        case 'Rain': return 'visibility-rain-overlay'
        default: return ''
    }
})
</script>