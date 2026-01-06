<template>
    <div class="relative inline-flex items-center justify-center" :style="{ width: `${size}px`, height: `${size}px` }">
        <!-- Background circle -->
        <svg class="absolute inset-0 -rotate-90" :width="size" :height="size">
            <circle
                :cx="size / 2"
                :cy="size / 2"
                :r="radius"
                fill="none"
                stroke="currentColor"
                class="text-[hsl(var(--maz-border))]"
                :stroke-width="strokeWidth"
            />
            <!-- Progress circle -->
            <circle
                :cx="size / 2"
                :cy="size / 2"
                :r="radius"
                fill="none"
                :stroke="progressColor"
                :stroke-width="strokeWidth"
                stroke-linecap="round"
                :stroke-dasharray="circumference"
                :stroke-dashoffset="strokeDashoffset"
                class="transition-all duration-300"
            />
        </svg>
        <!-- Time display -->
        <span class="text-lg font-bold font-mono tabular-nums z-10" :class="textClass">
            {{ formattedTime }}
        </span>
    </div>
</template>

<script setup lang="ts">
interface Props {
    duration?: number
    loop?: boolean
    autoStart?: boolean
    size?: number
    strokeWidth?: number
}

const props = withDefaults(defineProps<Props>(), {
    duration: 30,
    loop: false,
    autoStart: true,
    size: 72,
    strokeWidth: 4
})

const emit = defineEmits<{
    complete: []
    tick: [remaining: number]
}>()

const remaining = ref(props.duration)
const isRunning = ref(false)
let intervalId: ReturnType<typeof setInterval> | null = null

const radius = computed(() => (props.size - props.strokeWidth) / 2)
const circumference = computed(() => 2 * Math.PI * radius.value)
const progress = computed(() => remaining.value / props.duration)
const strokeDashoffset = computed(() => circumference.value * (1 - progress.value))

const formattedTime = computed(() => {
    const mins = Math.floor(remaining.value / 60)
    const secs = remaining.value % 60
    return mins > 0 
        ? `${mins}:${secs.toString().padStart(2, '0')}`
        : `${secs}`
})

const progressColor = computed(() => {
    if (remaining.value <= 5) return 'hsl(var(--maz-danger))'
    if (remaining.value <= 10) return 'hsl(var(--maz-warning))'
    return 'hsl(var(--maz-primary))'
})

const textClass = computed(() => {
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
