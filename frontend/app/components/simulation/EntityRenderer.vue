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
                parentHover && entities.length > 1 ? 'pointer-events-auto cursor-pointer hover:scale-125 z-50' : '',
                rankStore.isSelected(entity.id) ? 'selected-glow scale-110' : ''
            ]"
            :title="entity.metadata.name"
            @click.stop="handleEntityClick(entity)"
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

function handleEntityClick(entity: Entity) {
    if (props.entities.length === 1) return // Let parent handle single entity
    
    if (rankStore.isSelected(entity.id)) {
        rankStore.clearSelection()
    } else {
        rankStore.select(entity)
    }
}

function getStackStyles(index: number, total: number) {
    const zIndex = (index + 1) * 10
    
    if (total === 1) return { zIndex, transform: 'scale(1)' }

    if (props.parentHover) {
        const spacing = 28
        const offset = (index * spacing) - ((total - 1) * spacing / 2)
        return {
            zIndex: 100 + index,
            transform: `translateX(${offset}px) scale(1.15)`
        }
    }

    const offset = (index * 3) - ((total - 1) * 1.5)
    return {
        zIndex,
        transform: `translate(${offset}px, ${offset}px)`
    }
}
</script>

<style scoped>
.selected-glow {
    background: radial-gradient(circle, hsl(var(--maz-primary) / 0.4) 0%, transparent 70%);
    border-radius: 50%;
    padding: 0.25rem;
    animation: pulse-glow 1.5s ease-in-out infinite;
    filter: drop-shadow(0 0 8px hsl(var(--maz-primary) / 0.8));
}

@keyframes pulse-glow {
    0%, 100% {
        filter: drop-shadow(0 0 10px hsl(var(--maz-primary) / 0.8));
        background: radial-gradient(circle, hsl(var(--maz-primary) / 0.4) 0%, transparent 70%);
    }
    50% {
        filter: drop-shadow(0 0 18px hsl(var(--maz-primary) / 1));
        background: radial-gradient(circle, hsl(var(--maz-primary) / 0.6) 0%, transparent 70%);
    }
}
</style>