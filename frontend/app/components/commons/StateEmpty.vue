<template>
  <div class="flex flex-col items-center justify-center py-16 px-4">
    <!-- Icon container -->
    <div 
      class="w-20 h-20 rounded-2xl flex items-center justify-center mb-6"
      :class="iconContainerClass"
    >
      <component 
        :is="icon" 
        class="w-10 h-10"
        :class="iconClass"
      />
    </div>

    <!-- Title -->
    <h3 class="text-xl font-semibold text-center mb-2">
      {{ title }}
    </h3>

    <!-- Description -->
    <p class="text-[hsl(var(--maz-foreground))]/60 text-center max-w-sm mb-6">
      {{ description }}
    </p>

    <!-- Action button -->
    <MazBtn 
      v-if="actionLabel" 
      :color="actionColor"
      @click="$emit('action')"
    >
      {{ actionLabel }}
    </MazBtn>
  </div>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'
import { MazInboxStack } from '@maz-ui/icons'

const props = withDefaults(defineProps<{
  icon?: Component
  title?: string
  description?: string
  actionLabel?: string
  actionColor?: 'primary' | 'secondary' | 'info' | 'warning' | 'danger' | 'success'
  variant?: 'default' | 'warning' | 'info'
}>(), {
  icon: () => MazInboxStack,
  title: 'No Data Available',
  description: 'There is nothing to display at the moment.',
  actionColor: 'primary',
  variant: 'default'
})

defineEmits<{
  action: []
}>()

const iconContainerClass = computed(() => {
  switch (props.variant) {
    case 'warning': return 'bg-[hsl(var(--maz-warning))]/10'
    case 'info': return 'bg-[hsl(var(--maz-info))]/10'
    default: return 'bg-[hsl(var(--maz-primary))]/10'
  }
})

const iconClass = computed(() => {
  switch (props.variant) {
    case 'warning': return 'text-[hsl(var(--maz-warning))]'
    case 'info': return 'text-[hsl(var(--maz-info))]'
    default: return 'text-[hsl(var(--maz-primary))]'
  }
})
</script>
