<template>
    <div 
        class="relative w-full h-full flex items-center justify-center pointer-events-none"
    >
        <span 
            v-for="(entity, index) in entities"
            :key="entity.id"
            class="absolute transition-all duration-300 select-none text-base sm:text-lg md:text-xl lg:text-3xl"
            :style="getStackStyles(index, entities.length)"
            :class="[
                entity.metadata.is_occluded && !parentHover ? 'opacity-50 grayscale' : '',
                parentHover ? 'z-50' : '',
                rankStore.isSelected(entity.id) ? 'ring-2 ring-cyan-400 rounded-full bg-cyan-400/20 scale-110' : ''
            ]"
            :title="entity.metadata.name"
        >
            {{ entity.emoji }}
        </span>

        <div v-if="entities.length > 1 && !parentHover" class="absolute -top-1 -right-1 bg-red-700 text-[10px] text-white px-1 rounded-full z-40">
            {{ entities.length }}
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Entity } from '~/types/simulation'
import { useRankStore } from '~/stores/rank'

const props = defineProps<{
    entities: Entity[]
    parentHover: boolean
    hasSelected?: boolean
}>()

const rankStore = useRankStore()

function getStackStyles(index: number, total: number) {
    const zIndex = (index + 1) * 10
    
    if (total === 1) return { zIndex, transform: 'scale(1)' }

    if (props.parentHover) {
        const offset = (index * 20) - ((total - 1) * 6)
        return {
            zIndex: 100 + index,
            transform: `translateX(${offset}px) scale(1.1)`
        }
    }

    const offset = (index * 3) - ((total - 1) * 1.5)
    return {
        zIndex,
        transform: `translate(${offset}px, ${offset}px)`
    }
}
</script>