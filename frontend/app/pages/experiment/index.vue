<template>
    <div>
        <!-- Loading State -->
        <CommonsStateLoaderWheels v-if="store.isLoading" :text="loadingMessage" />

        <!-- Scenario Content -->
        <ExperimentScenarioViewer
            v-else-if="scenario"
            :scenario="scenario"
            :ranking-options="rankingOptions"
            :highlighted-zone="highlightedZone"
            :is-loading="store.isLoading"
            :show-timer="true"
            :timer-duration="timerDuration"
            :initial-time="computedInitialTime"
            :auto-start-timer="true"
            :header-title="`Scenario ${scenario.current_step} of ${scenario.total_steps}`"
            ranking-title="Which action are you likely to do? Rank these?"
            submit-button-text="Submit Response"
            @update:ranking-options="rankingOptions = $event"
            @highlight="highlightedZone = $event"
            @interaction="hasUserInteracted = $event"
            @submit="submitResponse(false)"
            @time-up="handleTimeUp"
        >
            <template #header-actions>
                <MazBtn 
                    size="sm"
                    color="destructive" 
                    outlined
                    :disabled="store.isLoading"
                    @click="showQuitDialog = true"
                >
                    <MazXCircle class="w-6 h-6" />
                    Quit
                </MazBtn>
            </template>
        </ExperimentScenarioViewer>


        <!-- Error/Empty State -->
        <CommonsStateError 
            v-else
            title="Failed to Load Scenario"
            description="We couldn't load the scenario. Please try again."
            :show-home="true"
            @retry="loadScenario"
        />

        <!-- Quit Confirmation Dialog -->
        <MazDialog v-model="showQuitDialog">
            <template #title>Quit Experiment?</template>
            <template #default>
                <p class="text-[hsl(var(--maz-muted))]">Your progress will be lost and you'll need to start over. Are you sure?</p>
            </template>
            <template #footer="{ close }">
                <div class="flex gap-3 justify-end">
                    <MazBtn color="primary" outlined @click="close">Continue</MazBtn>
                    <MazBtn color="destructive" outlined @click="handleQuit">Quit</MazBtn>
                </div>
            </template>
        </MazDialog>

        
    </div>
</template>

<script setup lang="ts">
import { useExperimentStore } from '~/stores/experiment'
import type { ScenarioResponse } from '~/types/response.types'
import { MazXCircle } from '@maz-ui/icons'

definePageMeta({
    layout: 'optional-layout',
    middleware: ['session']
})

const config = useRuntimeConfig()
const store = useExperimentStore()
const router = useRouter()

// Timer configuration
const timerDuration = computed(() => Number(config.public.timerDuration) || 20)
const startTime = ref<number>(0)

// State
const scenario = ref<ScenarioResponse | null>(null)
const rankingOptions = ref<{ key: string; label: string; zone: string }[]>([])
const highlightedZone = ref<string | null>(null)
const showQuitDialog = ref(false)
const hasUserInteracted = ref(false)
const loadingState = ref<'loading' | 'submitting'>('loading')
const maxTimeoutId = ref<ReturnType<typeof setTimeout> | null>(null)
const computedInitialTime = ref(20)

// Loading message based on state
const loadingMessage = computed(() => {
    return loadingState.value === 'submitting' 
        ? 'Submitting response...' 
        : 'Loading scenario...'
})

// Load scenario
async function loadScenario() {
    loadingState.value = 'loading'
    
    // Mark guide as completed when experiment starts
    store.completeGuide()
    
    // Clear any existing timeout
    if (maxTimeoutId.value) {
        clearTimeout(maxTimeoutId.value)
        maxTimeoutId.value = null
    }
    
    const data = await store.getScenario()
    if (data) {
        // Calculate elapsed time using localStorage (prevents refresh exploit)
        const storageKey = `scenario_start_${data.id}`
        let startedAt = localStorage.getItem(storageKey)
        
        if (!startedAt) {
            startedAt = Date.now().toString()
            localStorage.setItem(storageKey, startedAt)
        }
        
        const elapsedSeconds = Math.floor((Date.now() - parseInt(startedAt)) / 1000)
        let remaining = timerDuration.value - elapsedSeconds
        if (remaining < 0) remaining = 0
        
        computedInitialTime.value = remaining
        
        // Set start time for response calculation
        startTime.value = parseInt(startedAt)
        
        scenario.value = data
        // Reset interaction tracking for new scenario
        hasUserInteracted.value = false
        
        // Randomize ranking options to prevent passive agreement bias
        const options = [
            { key: 'maintain', label: data.dilemma_options.maintain, zone: 'zone_a' },
            { key: 'swerve_left', label: data.dilemma_options.swerve_left, zone: 'zone_b' },
            { key: 'swerve_right', label: data.dilemma_options.swerve_right, zone: 'zone_c' },
        ]
        for (let i = options.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [options[i], options[j]] = [options[j]!, options[i]!]
        }
        rankingOptions.value = options
        
        // Set backup timeout (timer duration + 3s buffer)
        maxTimeoutId.value = setTimeout(() => {
            handleTimeUp()
        }, (timerDuration.value + 3) * 1000)
        
        if (remaining === 0) {
            handleTimeUp()
        }
    }
}

async function handleTimeUp() {
    if (store.isLoading) return // Prevent double submission
    
    const responseTimeMs = Date.now() - startTime.value
    await submitResponse(true, responseTimeMs)
}

async function submitResponse(isTimeout = false, responseTimeMs?: number) {
    if (store.isLoading || !scenario.value) return
    
    // Clear backup timeout
    if (maxTimeoutId.value) {
        clearTimeout(maxTimeoutId.value)
        maxTimeoutId.value = null
    }
    
    loadingState.value = 'submitting'
    
    try {
        const finalResponseTime = responseTimeMs || (Date.now() - startTime.value)
        
        await store.submitResponse(scenario.value.id, {
            ranking_order: rankingOptions.value.map(opt => opt.key),
            response_time_ms: finalResponseTime,
            is_timeout: isTimeout,
            has_interacted: hasUserInteracted.value
        })
        
        // Clean up localStorage for this scenario
        localStorage.removeItem(`scenario_start_${scenario.value.id}`)
        
        // Check if experiment is complete
        if (scenario.value.current_step >= scenario.value.total_steps) {
            // Keep loading visible during redirect
            store.isLoading = true
            await new Promise(resolve => setTimeout(resolve, 1000))
            
            // Redirect to thank you page (token cleared after feedback fetch)
            await router.push('/thank-you')
            return
        }
        
        // Load next scenario
        await loadScenario()
        
    } catch (error) {
        // Error handling is done in the store
        loadingState.value = 'loading'
    }
}

function handleQuit() {
    store.token = null
    router.push('/')
}

onMounted(loadScenario)

onUnmounted(() => {
    // Clean up timeout
    if (maxTimeoutId.value) {
        clearTimeout(maxTimeoutId.value)
        maxTimeoutId.value = null
    }
})
</script>
