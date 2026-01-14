<template>
    <div class="flex flex-col gap-3">
        <VueDraggableNext 
            :list="localOptions" 
            :animation="200"
            ghost-class="opacity-50"
            @start="onDragStart"
            @end="onDragEnd"
            class="flex flex-col gap-3"
        >
            <div 
                v-for="(element, index) in localOptions"
                :key="element.key"
                class="flex items-center gap-4 p-4 rounded-xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-secondary))] transition-colors cursor-grab active:cursor-grabbing hover:border-[hsl(var(--maz-primary))]"
                @mouseenter="onHover(element.zone)"
                @mouseleave="onHover(null)"
            >
                <span class="flex items-center justify-center w-8 h-8 rounded-full bg-[hsl(var(--maz-primary))]/20 text-[hsl(var(--maz-primary))] font-bold text-sm">
                    {{ index + 1 }}
                </span>
                <div class="flex-1">
                    <p class="font-medium">{{ element.label }}</p>
                </div>
                <MazBars3 class="w-5 h-5 text-[hsl(var(--maz-muted))]" />
            </div>
        </VueDraggableNext>
    </div>
</template>

<script setup lang="ts">
import { VueDraggableNext } from 'vue-draggable-next'
import { ref, watch } from 'vue'
import { MazBars3 } from '@maz-ui/icons'

export interface RankingOption {
    key: string
    label: string
    zone: string
}

const props = defineProps<{
    options: RankingOption[]
}>()

const emit = defineEmits<{
    'update:options': [options: RankingOption[]]
    'highlight': [zone: string | null]
    'interaction': [hasInteracted: boolean]
}>()

const dragging = ref(false)
const hasInteracted = ref(false)
const interactionCount = ref(0)

const localOptions = ref([...props.options])

// Watch for external changes - only sync if arrays are actually different
watch(() => props.options, (newVal) => {
    const newKeys = newVal.map(o => o.key).join(',')
    const localKeys = localOptions.value.map(o => o.key).join(',')
    if (newKeys !== localKeys) {
        localOptions.value = [...newVal]
    }
}, { deep: true })

function trackInteraction() {
    if (hasInteracted.value) return
    interactionCount.value++
    if (interactionCount.value >= 2) {
        hasInteracted.value = true
        emit('interaction', true)
    }
}

function onDragStart() {
    dragging.value = true
    trackInteraction()
}

function onDragEnd() {
    dragging.value = false
    // Emit the final order after drag completes
    emit('update:options', [...localOptions.value])
    emit('highlight', null)
}

function onHover(zone: string | null) {
    if (dragging.value) return
    emit('highlight', zone)
    if (zone !== null) {
        trackInteraction()
    }
}
</script>
