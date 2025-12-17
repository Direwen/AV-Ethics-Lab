<template>
    <div 
        class="relative w-full h-full flex items-center justify-center pointer-events-none"
    >
        <span 
            v-for="(entity, index) in entities"
            :key="entity.id"
            class="absolute transition-all duration-300 select-none text-base sm:text-lg md:text-xl lg:text-3xl pointer-events-auto"
            :style="getStackStyles(index, entities.length)"
            :class="{ 
                'opacity-50 grayscale': entity.metadata.is_occluded && !isHovered,
                'z-50': isHovered // Pop to top on hover
            }"
            :title="entity.metadata.name"
            @mouseenter="isHovered = true"
            @mouseleave="isHovered = false"
        >
            {{ entity.emoji }}
        </span>

        <div v-if="entities.length > 1 && !isHovered" class="absolute -top-1 -right-1 bg-red-700 text-[10px] text-white px-1 rounded-full z-40">
            {{ entities.length }}
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Entity } from '~/types/simulation'

const props = defineProps<{
    entities: Entity[]
}>()

const isHovered = ref(false)

// --- STACKING LOGIC ---
function getStackStyles(index: number, total: number) {
    // Base Z-Index to ensure proper layering
    const zIndex = (index + 1) * 10
    
    if (total === 1) return { zIndex, transform: 'scale(1)' }

    // If hovered, "explode" the stack so user can see everything
    if (isHovered.value) {
        // Spread horizontally
        const offset = (index * 12) - ((total - 1) * 6)
        return {
            zIndex: 100 + index, // High Z-Index for expanded view
            transform: `translateX(${offset}px) scale(1.1)`
        }
    }

    // Default: Tight stack
    const offset = (index * 3) - ((total - 1) * 1.5)
    return {
        zIndex,
        transform: `translate(${offset}px, ${offset}px)`
    }
}
</script>