<template>
    <div>
        <MazBtn @click="toggleDarkMode">
            Toggle dark mode ({{ isDark ? '🌙' : '☀️' }})
        </MazBtn>
        <MazBtn @click="toast.error('whahah')">
            Toast
        </MazBtn>
        <MazBtn @click="testCreateSession" :loading="experimentStore.isLoading">
            Test Create Session
        </MazBtn>
    </div>
</template>

<script setup>
const { toggleDarkMode, isDark } = useMazTheme()
const toast = useMazToast()
const experimentStore = useExperimentStore()

async function testCreateSession() {
    await experimentStore.createSession({}, true)
}
</script>