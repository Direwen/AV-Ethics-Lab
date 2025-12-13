<template>
    <div 
        class="flex items-center justify-center transition-all duration-300 ease-out 
                custom-background-gradient
                text-[hsl(var(--maz-foreground))]
                cursor-pointer 
                shadow-lg shadow-primary/30 hover:shadow-primary/50 hover:-translate-y-0.5"
        :class="containerSize[size]"
    >
        <MazSquares2x2 
        class="text-white drop-shadow-sm" 
        :class="iconSize[size]"
        />
    </div>
</template>

<script setup lang="ts">
import { MazSquares2x2 } from '@maz-ui/icons'

interface Props {
    size?: 'sm' | 'md' | 'lg' | 'xl'
}

withDefaults(defineProps<Props>(), {
    size: 'md'
})

const containerSize = {
    sm: 'w-8 h-8 rounded-lg',
    md: 'w-10 h-10 rounded-xl',
    lg: 'w-14 h-14 rounded-2xl',
    xl: 'w-20 h-20 rounded-3xl'
}

const iconSize = {
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-7 h-7',
    xl: 'w-10 h-10'
}
</script>