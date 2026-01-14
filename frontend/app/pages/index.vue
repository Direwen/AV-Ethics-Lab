<template>
    <div
        class="flex flex-col items-center justify-center w-full max-w-5xl mx-auto space-y-12 sm:space-y-16 px-4"
    >
        <section class="text-center space-y-6 sm:space-y-8 py-10 sm:py-12 relative z-10">
            <div
                class="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-[hsl(var(--maz-primary))] bg-[hsl(var(--maz-primary))]/10 text-sm font-medium text-[hsl(var(--maz-primary))]"
            >
                <span class="relative flex h-2 w-2">
                    <span
                        class="animate-ping absolute inline-flex h-full w-full rounded-full bg-[hsl(var(--maz-primary))] opacity-75"
                    ></span>
                    <span
                        class="relative inline-flex rounded-full h-2 w-2 bg-[hsl(var(--maz-primary))]"
                    ></span>
                </span>
                Active Research Study
            </div>

            <h1 class="text-3xl sm:text-4xl md:text-6xl font-bold tracking-tight leading-tight">
                <span class="whitespace-nowrap">Ethical&nbsp;Decisions</span> in
                <span
                    class="text-transparent bg-clip-text bg-gradient-to-br from-[hsl(var(--maz-primary))] to-[hsl(var(--maz-accent))] whitespace-nowrap"
                >
                    Autonomous&nbsp;Vehicles
                </span>
            </h1>

            <p
                class="text-base sm:text-lg md:text-xl text-[hsl(var(--maz-foreground))] opacity-70 max-w-2xl mx-auto leading-relaxed"
            >
                Contribute to academic research on moral decision-making in autonomous vehicles. Your responses help shape ethical frameworks for self-driving technology.
            </p>

            <div
                class="flex items-center justify-center gap-3 sm:gap-4 pt-4"
            >
                <NuxtLink to="/experiment/consent">
                    <MazBtn
                        size="lg"
                        class="md:size-xl rounded-xl shadow-xl shadow-[hsl(var(--maz-primary))]/20 custom-background-gradient hover:scale-105 transition-transform"
                    >
                        Start Experiment
                    </MazBtn>
                </NuxtLink>

                <NuxtLink to="/dashboard">
                    <MazBtn
                        pastel
                        color="background"
                        variant="text"
                        size="lg"
                        class="md:size-xl shadow hover:scale-105 transition-transform"
                    >
                        Dashboard
                    </MazBtn>
                </NuxtLink>
            </div>
        </section>

        <section class="grid grid-cols-1 md:grid-cols-3 gap-6 w-full relative z-10">
            <div
                v-for="(feature, index) in features"
                :key="index"
                class="p-6 rounded-2xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-background))]/50 backdrop-blur-sm hover:border-[hsl(var(--maz-primary))] transition-colors group"
            >
                <div
                    class="w-12 h-12 rounded-lg bg-[hsl(var(--maz-primary))]/10 flex items-center justify-center mb-4 text-[hsl(var(--maz-primary))] group-hover:scale-110 transition-transform"
                >
                    <component :is="feature.icon" class="h-6 w-6" />
                </div>
                <h3 class="text-lg font-semibold mb-2">{{ feature.title }}</h3>
                <p class="text-sm text-[hsl(var(--maz-foreground))] opacity-60">
                    {{ feature.description }}
                </p>
            </div>
        </section>

        <div
            class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full h-full max-w-4xl max-h-96 -z-0 opacity-20 pointer-events-none"
        >
            <div class="grid grid-cols-12 gap-4 h-full mask-linear-fade">
                <div
                    v-for="n in 48"
                    :key="n"
                    class="custom-background-gradient rounded-sm animate-pulse"
                    :style="{
                        opacity: Math.random() * 0.5,
                        animationDelay: `${Math.random() * 2}s`,
                    }"
                ></div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { MazEyeSlash, MazClock, MazCpuChip } from '@maz-ui/icons'

definePageMeta({
    middleware: "session",
});

const features = [
    {
        icon: MazEyeSlash,
        title: "100% Anonymous",
        description:
            "No personal data, IP addresses, or tracking cookies are stored. Your decisions cannot be traced back to you.",
    },
    {
        icon: MazClock,
        title: "5-10 Minutes",
        description:
            "A focused simulation designed to capture your intuitive moral judgments in realistic driving scenarios.",
    },
    {
        icon: MazCpuChip,
        title: "Real-World Impact",
        description:
            'Data from this study contributes directly to academic research on ethical decision-making in autonomous vehicles.',
    },
];
</script>

<style scoped>
/* Fade out the grid decoration at the edges */
.mask-linear-fade {
    mask-image: radial-gradient(circle at center, black 40%, transparent 70%);
}
</style>
