<template>
    <div class="flex flex-col min-h-screen">

        <header 
            class="sticky top-0 z-50 w-full border-b border-[hsl(var(--maz-border))] backdrop-blur-md transition-all duration-300"
        >
            <div class="container mx-auto flex h-16 items-center justify-between px-4 md:px-6">
            
            <NuxtLink to="/" class="flex items-center gap-3 group">
                <Logo size="md" />
                <span class="text-lg font-semibold tracking-tight group-hover:text-[hsl(var(--maz-primary))] transition-colors">
                AV Ethics Lab
                </span>
            </NuxtLink>

            <nav class="hidden md:flex items-center gap-6">
                
                <div class="flex items-center gap-2">
                    <ThemeToggle />
                </div>
                
                <NuxtLink to="/experiment/consent">
                <MazBtn 
                    size="sm" 
                    class="custom-background-gradient rounded-lg shadow-lg shadow-primary/20 hover:shadow-primary/40 transition-all"
                >
                    Get Started
                </MazBtn>
                </NuxtLink>
            </nav>

            </div>
        </header>

        <main class="flex-grow flex flex-col container mx-auto px-4 md:px-6 py-8">
            <slot />
        </main>

        <footer class="border-t border-[hsl(var(--maz-border))] py-8 mt-auto">
            <div
                class="container mx-auto px-4 md:px-6 flex flex-col md:flex-row justify-between items-center gap-4 text-sm text-[hsl(var(--maz-muted))]"
            >
                <div class="flex flex-col md:flex-row items-center gap-2 text-center md:text-left">
                    <span class="opacity-70">
                        &copy; {{ new Date().getFullYear() }} Min Sik Hein
                    </span>
                    <span class="hidden md:inline opacity-30">|</span>
                    <span class="opacity-70">University of Sunderland</span>
                </div>

                <div class="flex flex-col md:flex-row items-center gap-4 md:gap-6">
                    <a
                        href="mailto:myomyothannaing.1983@gmail.com"
                        class="hover:text-[hsl(var(--maz-primary))] transition-colors"
                    >
                        Supervisor: Dr. Myo Myo Than Naing
                    </a>
                    
                    <span class="hidden md:inline opacity-30">|</span>
                    
                    <span class="opacity-70" title="Computing Project">
                        Module: CET300
                    </span>
                </div>
            </div>
        </footer>

    </div>
</template>

<script setup lang="ts">
import Logo from '~/components/commons/Logo.vue';
import ThemeToggle from '~/components/commons/ThemeToggle.vue';
</script>