<template>
    <div ref="wrapperRef" class="w-full flex justify-center">
        <div 
            class="relative inline-block origin-top transition-transform duration-150"
            :style="{ transform: `scale(${scale})` }"
        >
            <!-- Main content with visibility filter -->
            <div 
                ref="contentRef"
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
        <!-- Spacer to maintain layout height -->
        <div :style="{ height: `${scaledHeight}px`, marginTop: `-${scaledHeight}px` }" class="pointer-events-none" />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'

export type Visibility = 'Clear' | 'Fog' | 'Night' | 'Rain'

const props = defineProps<{
    visibility?: Visibility
}>()

const wrapperRef = ref<HTMLElement | null>(null)
const contentRef = ref<HTMLElement | null>(null)
const containerWidth = ref(0)
const contentWidth = ref(0)
const contentHeight = ref(0)

// Scale factor: fit content to container, max 1.0
const scale = computed(() => {
    if (!contentWidth.value || !containerWidth.value) return 1
    return Math.min(1, containerWidth.value / contentWidth.value)
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
        contentWidth.value = contentRef.value.scrollWidth
        contentHeight.value = contentRef.value.scrollHeight
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
