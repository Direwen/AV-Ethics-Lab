<template>
    <div class="text-center max-w-2xl bg-red-500">
        <h2 class="text-4xl font-bold mb-6">Autonomous Vehicle Ethics Study</h2>
        <p class="mb-8 text-lg text-slate-400">
            Participate in a 5-minute interactive simulation to help researchers understand moral decision-making in AI.
        </p>
        
        <NuxtLink to="/auth/consent">
            <button class="bg-sensor-gradient text-white px-8 py-3 rounded-lg font-semibold hover:shadow-cyan-500/20 transition-all">
                Participate Now
            </button>
        </NuxtLink>
        <MazBtn @click="toggleDarkMode">
            Toggle dark mode ({{ isDark ? '🌙' : '☀️' }})
        </MazBtn>
    </div>
</template>

<script setup type="ts">

definePageMeta({
    middleware: 'session'
})

const { toggleDarkMode, isDark } = useMazTheme()
const toast = useMazToast()
const experimentStore = useExperimentStore()

async function testCreateSession() {
    await experimentStore.createSession({}, true)
}
</script>