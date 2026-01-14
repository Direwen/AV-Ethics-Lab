<template>
    <div class="flex flex-col min-h-screen bg-gradient-to-br from-[hsl(var(--maz-background))] to-[hsl(var(--maz-background-accent))]">
        
        <header class="w-full border-b border-[hsl(var(--maz-border))] backdrop-blur-md transition-all duration-300 relative z-20">
            <div class="container mx-auto flex h-16 items-center justify-between px-4 md:px-6">
            
                <div class="flex items-center gap-4">
                    <button 
                        @click="$router.back()" 
                        class="p-2 rounded-lg hover:bg-[hsl(var(--maz-secondary))] text-[hsl(var(--maz-foreground))] transition-colors"
                        aria-label="Go Back"
                    >
                        <MazArrowSmallLeft class="w-6 h-6" />
                    </button>
                    
                    <div class="h-6 w-px bg-[hsl(var(--maz-border))] mx-2 hidden sm:block"></div>

                    <NuxtLink to="/" class="flex items-center gap-3 group">
                        <Logo size="sm" />
                        <span class="text-lg font-semibold tracking-tight group-hover:text-[hsl(var(--maz-primary))] transition-colors">
                            AV Ethics Lab
                        </span>
                    </NuxtLink>
                </div>

                <div class="flex items-center gap-4">
                    <CommonThemeToggle />
                </div>

            </div>
        </header>

        <main class="flex-grow flex flex-col justify-center py-8">
            <slot />
        </main>

    </div>
</template>

<script setup lang="ts">
import { MazArrowSmallLeft } from '@maz-ui/icons'
</script>