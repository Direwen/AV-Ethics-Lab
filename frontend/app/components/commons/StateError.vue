<template>
  <div class="flex flex-col items-center justify-center py-16 px-4">
    <!-- Error icon with pulse ring -->
    <div class="relative mb-6">
      <div class="absolute inset-0 w-20 h-20 rounded-2xl bg-[hsl(var(--maz-danger))]/20 animate-pulse" />
      <div class="relative w-20 h-20 rounded-2xl bg-[hsl(var(--maz-danger))]/10 flex items-center justify-center">
        <MazExclamationTriangle class="w-10 h-10 text-[hsl(var(--maz-danger))]" />
      </div>
    </div>

    <!-- Title -->
    <h3 class="text-xl font-semibold text-center mb-2">
      {{ title }}
    </h3>

    <!-- Description -->
    <p class="text-[hsl(var(--maz-foreground))]/60 text-center max-w-sm mb-2">
      {{ description }}
    </p>

    <!-- Error details (collapsible) -->
    <details v-if="errorDetails" class="mb-6 w-full max-w-sm">
      <summary class="text-sm text-[hsl(var(--maz-foreground))]/40 cursor-pointer hover:text-[hsl(var(--maz-foreground))]/60 transition-colors">
        Show details
      </summary>
      <pre class="mt-2 p-3 rounded-lg bg-[hsl(var(--maz-danger))]/5 text-xs text-[hsl(var(--maz-danger))] overflow-auto">{{ errorDetails }}</pre>
    </details>

    <!-- Actions -->
    <div class="flex items-center gap-3">
      <MazBtn 
        v-if="showRetry"
        color="primary"
        @click="$emit('retry')"
      >
        <MazArrowPath class="w-4 h-4 mr-2" />
        Try Again
      </MazBtn>
      
      <NuxtLink v-if="showHome" to="/">
        <MazBtn 
          color="secondary"
          variant="outline"
        >
        Go Home
      </MazBtn>
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { MazExclamationTriangle, MazArrowPath } from '@maz-ui/icons'

withDefaults(defineProps<{
  title?: string
  description?: string
  errorDetails?: string
  showRetry?: boolean
  showHome?: boolean
}>(), {
  title: 'Something Went Wrong',
  description: 'We encountered an unexpected error. Please try again.',
  showRetry: true,
  showHome: false
})

defineEmits<{
  retry: []
}>()
</script>
