<template>
    <div 
        class="relative w-full h-full flex items-center justify-center pointer-events-none"
    >
        <span 
            v-for="(entity, index) in entities"
            
            :key="entity.id"
            class="absolute transition-all duration-300 select-none text-base sm:text-lg md:text-xl lg:text-3xl"
            :style="getEntityStyles(entity, index, entities.length)"
            :class="[
                parentHover && entities.length > 1 ? 'pointer-events-auto cursor-pointer z-50' : '',
                rankStore.isSelected(entity.id) ? 'selected-glow' : ''
            ]"
            :title="`${entity.type} -> ${entity.metadata?.orientation}`"            @click.stop="handleEntityClick(entity)"
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

function getOrientationTransform(orientation: string | undefined): string {
    if (!orientation) return ''
    
    // All emojis assumed to face LEFT (West) by default
    switch (orientation) {
        case 'N': return 'rotate(90deg)'  // Left -> Up
        case 'S': return 'rotate(-90deg)' // Left -> Down
        case 'E': return 'scaleX(-1)'     // Left -> Flip to Right
        case 'W': return 'rotate(0deg)'   // Left -> Left (Natural)
        default: return ''
    }
}

function getEntityStyles(entity: Entity, index: number, total: number) {
    const zIndex = (index + 1) * 10
    const orientationTransform = getOrientationTransform(entity.metadata?.orientation)
    
    // Base transforms array
    let transformList: string[] = []

    // 1. Handle Stacking / Hover Offsets
    if (total > 1) {
        if (props.parentHover) {
            const spacing = 28
            const offset = (index * spacing) - ((total - 1) * spacing / 2)
            transformList.push(`translateX(${offset}px)`)
            transformList.push('scale(1.15)')
        } else {
            const offset = (index * 3) - ((total - 1) * 1.5)
            transformList.push(`translate(${offset}px, ${offset}px)`)
        }
    } else {
        // Even for single items, explicit scale(1) helps prevent fuzzy rendering text issues in some browsers
        transformList.push('scale(1)')
    }

    // 2. Add Orientation (Last, so it rotates IN PLACE after translation)
    if (orientationTransform) {
        transformList.push(orientationTransform)
    }

    return { 
        zIndex: props.parentHover ? 100 + index : zIndex, 
        transform: transformList.join(' ') 
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