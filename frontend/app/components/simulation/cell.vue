<template>
    <div
        ref="cellRef"
        class="relative w-4 h-4 sm:w-6 sm:h-6 md:w-8 md:h-8 lg:w-12 lg:h-12 flex items-center justify-center transition-all duration-30"
        :class="[
            definition?.class, 
            isInteractive ? 'hover:brightness-110 cursor-pointer' : ''
        ]"
        :title="definition?.name" 
        role="gridcell"
        @click="handleClick"
        @mouseenter="isHovered = true"  
        @mouseleave="isHovered = false"
    >
        <EntityRenderer 
            v-if="entities?.length" 
            :entities="entities" 
            :parent-hover="isHovered"
            :has-selected="hasSelectedEntity"
        />
    </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import EntityRenderer from '~/components/simulation/EntityRenderer.vue'
import type { Entity, CellDefinition } from '~/types/simulation'
import { useRankStore } from '~/stores/rank'

const props = defineProps<{
    cellCode: number | string
    definition: CellDefinition
    entities?: Entity[]
}>()

const rankStore = useRankStore()
const isHovered = ref(false)

const isInteractive = computed(() => props.definition?.isInteractive ?? false)

const hasSelectedEntity = computed(() => {
    if (!props.entities?.length) return false
    return props.entities.some(e => rankStore.isSelected(e.id))
})

function handleClick() {
    // Non-interactive cells clear selection
    if (!isInteractive.value) {
        rankStore.clearSelection()
        return
    }

    // Empty cell - clear selection
    if (!props.entities?.length) {
        rankStore.clearSelection()
        return
    }

    // Single entity - toggle selection
    if (props.entities.length === 1) {
        const entity = props.entities[0]!
        if (rankStore.isSelected(entity.id)) {
            rankStore.clearSelection()
        } else {
            rankStore.select(entity)
        }
        return
    }

    // Multiple entities - TODO: show picker menu
}
</script>