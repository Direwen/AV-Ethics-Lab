<template>
  <div class="inline-block border-2 border-black">
    <div 
      v-for="row in 8" 
      :key="row" 
      class="flex"
    >
      <div 
        v-for="col in 8" 
        :key="col" 
        class="w-16 h-16 flex items-center justify-center text-2xl border border-gray-400"
        :class="(row + col) % 2 === 0 ? 'bg-black text-white' : 'bg-white text-black'"
      >
        {{ getCellContent(row, col) }}
      </div>
    </div>
  </div>
</template>

<script setup>
const emojis = ['🎮', '🎯', '🎲', '🎪', '🎨', '🎭', '🎸', '🎺', '🎻', '🎹', '🚀', '🌟', '⭐', '💎', '🔥'];

const getCellContent = (row, col) => {
  // Use row and col to create a deterministic pattern
  const seed = row * 8 + col;
  // Only show emoji in ~40% of cells
  if (seed % 5 < 2) {
    return emojis[seed % emojis.length];
  }
  return '';
};
</script>