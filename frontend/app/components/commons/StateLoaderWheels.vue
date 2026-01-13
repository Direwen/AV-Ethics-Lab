<template>
  <div class="flex flex-col items-center justify-center gap-6 py-12">
    <!-- Two wheels animation -->
    <div class="flex items-end gap-8">
      <!-- Left wheel -->
      <div class="wheel-container animate-spin" style="animation-duration: 0.8s;">
        <div class="wheel-outer custom-background-gradient" />
        <div class="wheel-inner bg-[hsl(var(--maz-background))]">
          <div class="wheel-hub custom-background-gradient opacity-50" />
          <!-- Spokes -->
          <div class="spoke spoke-v custom-background-gradient opacity-30" />
          <div class="spoke spoke-h custom-background-gradient opacity-30" />
        </div>
      </div>

      <!-- Right wheel -->
      <div class="wheel-container animate-spin" style="animation-duration: 0.8s;">
        <div class="wheel-outer custom-background-gradient" />
        <div class="wheel-inner bg-[hsl(var(--maz-background))]">
          <div class="wheel-hub custom-background-gradient opacity-50" />
          <!-- Spokes -->
          <div class="spoke spoke-v custom-background-gradient opacity-30" />
          <div class="spoke spoke-h custom-background-gradient opacity-30" />
        </div>
      </div>
    </div>

    <!-- Road line animation -->
    <div class="w-32 h-1 bg-[hsl(var(--maz-border))] rounded-full overflow-hidden">
      <div class="h-full w-8 custom-background-gradient rounded-full animate-road" />
    </div>

    <!-- Text -->
    <p v-if="text" class="text-[hsl(var(--maz-foreground))]/60 text-sm">
      {{ text }}
    </p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  text?: string
}>()
</script>

<style scoped>
.wheel-container {
  position: relative;
  width: 48px;
  height: 48px;
}

.wheel-outer {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  mask: 
    linear-gradient(#fff 0 0) content-box,
    linear-gradient(#fff 0 0);
  mask-composite: exclude;
  padding: 4px;
}

.wheel-inner {
  position: absolute;
  inset: 4px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.wheel-hub {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.spoke {
  position: absolute;
}

.spoke-v {
  width: 2px;
  height: 32px;
}

.spoke-h {
  width: 32px;
  height: 2px;
}

@keyframes road {
  0% { transform: translateX(-2rem); }
  100% { transform: translateX(8rem); }
}

.animate-road {
  animation: road 1s linear infinite;
}
</style>
