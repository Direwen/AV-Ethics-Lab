<template>
  <div class="space-y-4">
    <div
      v-for="n in count"
      :key="n"
      class="rounded-xl bg-[hsl(var(--maz-foreground))]/5 overflow-hidden"
      :class="sizeClasses"
    >
      <div class="h-full w-full animate-shimmer bg-gradient-to-r from-transparent via-[hsl(var(--maz-foreground))]/10 to-transparent" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  count?: number
  size?: 'sm' | 'md' | 'lg'
}>(), {
  count: 3,
  size: 'md'
})

const sizeClasses = computed(() => {
  switch (props.size) {
    case 'sm': return 'h-10'
    case 'lg': return 'h-32'
    default: return 'h-20'
  }
})
</script>

<style scoped>
@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}
.animate-shimmer {
  animation: shimmer 1.5s infinite;
}
</style>
