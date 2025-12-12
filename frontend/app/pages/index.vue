<template>

    <div>
        <MazBtn @click="toggleDarkMode">
            Toggle dark mode ({{ isDark ? '🌙' : '☀️' }})
        </MazBtn>
        <MazBtn @click="toast.error('whahah')">
            Toast
        </MazBtn>
    </div>
</template>

<script setup>
const { toggleDarkMode, isDark } = useMazTheme()
const toast = useMazToast()
</script>