<template>
  <div
    class="min-h-screen bg-radial
          from-[hsl(var(--maz-background-accent))]
          to-[hsl(var(--maz-background))]
          text-[hsl(var(--maz-foreground))]"
  >
    <NuxtLayout>
      <NuxtRouteAnnouncer />
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<script setup type="ts">
</script>
