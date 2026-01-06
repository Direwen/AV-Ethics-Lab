<template>
    <span class="text-2xl font-bold font-mono tabular-nums" :class="timerClass">
        {{ formattedTime }}
    </span>
</template>

<script setup lang="ts">
interface Props {
    duration?: number // in seconds
    loop?: boolean
    autoStart?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    duration: 30,
    loop: false,
    autoStart: true
})

const emit = defineEmits<{
    complete: []
    tick: [remaining: number]
}>()

const remaining = ref(props.duration)
const isRunning = ref(false)
let intervalId: ReturnType<typeof setInterval> | null = null

const formattedTime = computed(() => {
    const mins = Math.floor(remaining.value / 60)
    const secs = remaining.value % 60
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
})

const timerClass = computed(() => {
    if (remaining.value <= 5) return 'text-[hsl(var(--maz-danger))] animate-pulse'
    if (remaining.value <= 10) return 'text-[hsl(var(--maz-warning))]'
    return 'text-[hsl(var(--maz-foreground))]'
})

function start() {
    if (isRunning.value) return
    isRunning.value = true
    intervalId = setInterval(() => {
        remaining.value--
        emit('tick', remaining.value)
        
        if (remaining.value <= 0) {
            emit('complete')
            if (props.loop) {
                remaining.value = props.duration
            } else {
                stop()
            }
        }
    }, 1000)
}

function stop() {
    isRunning.value = false
    if (intervalId) {
        clearInterval(intervalId)
        intervalId = null
    }
}

function reset() {
    stop()
    remaining.value = props.duration
}

onMounted(() => {
    if (props.autoStart) start()
})

onUnmounted(() => {
    stop()
})

defineExpose({ start, stop, reset, remaining })
</script>
