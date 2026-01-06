<template>
  <div class="flex gap-4 text-sm text-gray-400">
    <div v-if="directions.includes('W')" class="flex items-center gap-1">
      <span class="text-yellow-400">←</span> West
    </div>
    <div v-if="directions.includes('E')" class="flex items-center gap-1">
      <span class="text-blue-400">→</span> East
    </div>
    <div v-if="directions.includes('N')" class="flex items-center gap-1">
      <span class="text-green-400">↑</span> North
    </div>
    <div v-if="directions.includes('S')" class="flex items-center gap-1">
      <span class="text-red-400">↓</span> South
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  directions: string[]
}>()
</script>
