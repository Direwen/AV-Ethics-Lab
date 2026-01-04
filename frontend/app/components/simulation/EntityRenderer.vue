<template>
    <div class="relative w-full h-full flex items-center justify-center pointer-events-none">
        <span 
            v-for="(entity, index) in entities"
            :key="entity.id"
            class="absolute transition-all duration-300 select-none text-base sm:text-lg md:text-xl lg:text-3xl"
            :style="getEntityStyles(entity, index, entities.length)"
            :class="getEntityClass(entity)"
            :title="`${entity.type} -> ${entity.metadata?.orientation}`"
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

const props = defineProps<{
    entities: Entity[]
    parentHover: boolean
}>()

function getEntityClass(entity: Entity): string {
    if (entity.metadata?.is_ego) return 'ego-glow'
    if (entity.metadata?.is_star) return 'star-glow'
    return ''
}

function getOrientationTransform(orientation: string | undefined): string {
    if (!orientation) return ''
    
    switch (orientation) {
        case 'N': return 'rotate(90deg)'
        case 'S': return 'rotate(-90deg)'
        case 'E': return 'scaleX(-1)'
        case 'W': return 'rotate(0deg)'
        default: return ''
    }
}

function getEntityStyles(entity: Entity, index: number, total: number) {
    const zIndex = (index + 1) * 10
    const orientationTransform = getOrientationTransform(entity.metadata?.orientation)
    
    let transformList: string[] = []

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
        transformList.push('scale(1)')
    }

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
.ego-glow {
    border-radius: 50%;
    padding: 0.15rem;
    filter: drop-shadow(0 0 6px rgba(59, 130, 246, 0.9)) drop-shadow(0 0 12px rgba(34, 197, 94, 0.6));
    animation: ego-pulse 2s ease-in-out infinite;
}

@keyframes ego-pulse {
    0%, 100% {
        filter: drop-shadow(0 0 6px rgba(59, 130, 246, 0.9)) drop-shadow(0 0 12px rgba(34, 197, 94, 0.6));
    }
    50% {
        filter: drop-shadow(0 0 10px rgba(59, 130, 246, 1)) drop-shadow(0 0 18px rgba(34, 197, 94, 0.8));
    }
}

.star-glow {
    border-radius: 50%;
    padding: 0.15rem;
    filter: drop-shadow(0 0 8px rgba(239, 68, 68, 0.9)) drop-shadow(0 0 14px rgba(249, 115, 22, 0.7));
    animation: star-pulse 1s ease-in-out infinite;
}

@keyframes star-pulse {
    0%, 100% {
        filter: drop-shadow(0 0 8px rgba(239, 68, 68, 0.9)) drop-shadow(0 0 14px rgba(249, 115, 22, 0.7));
    }
    50% {
        filter: drop-shadow(0 0 14px rgba(239, 68, 68, 1)) drop-shadow(0 0 22px rgba(249, 115, 22, 0.9));
    }
}
</style>
