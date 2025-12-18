<template>
    <MazSwitch
        v-model="isDarkMode"
        name="theme-toggle"
        color="primary"
        @update:model-value="toggleTheme"
    />

</template>

<script setup lang="ts">
import { useTheme } from '@maz-ui/themes';

const { colorMode, setColorMode } = useTheme();
const isDarkMode = ref(colorMode.value === 'dark');

const toggleTheme = (value: boolean) => {
    setColorMode(value ? 'dark' : 'light');
};
</script>