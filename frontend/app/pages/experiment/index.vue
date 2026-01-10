<template>
    <div class="container mx-auto px-4 py-6">
        <!-- Loading State -->
        <div v-if="store.isLoading" class="flex flex-col items-center justify-center min-h-[50vh]">
            <MazSpinner size="3em" />
            <p class="mt-4 text-sm text-[hsl(var(--maz-muted))]">{{ loadingMessage }}</p>
        </div>

        <!-- Scenario Content -->
        <template v-else-if="scenario">
            <!-- Header with Timer -->
            <div class="flex items-center justify-between mb-6">
                <div>
                    <h1 class="text-xl font-semibold">Scenario {{ scenario.current_step }} of {{ scenario.total_steps }}</h1>
                    <p class="text-sm text-[hsl(var(--maz-muted))]">{{ scenario.template_name }}</p>
                </div>
                <div class="flex items-center gap-4">
                    <ExperimentTimer 
                        v-if="scenario"
                        :key="scenario.id"
                        ref="timerRef"
                        :duration="timerDuration" 
                        :initial-time="computedInitialTime"
                        :auto-start="true"
                        @complete="handleTimeUp"
                    />
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
                </div>
            </div>

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

            <!-- Narrative -->
            <div class="mb-6 p-5 rounded-xl border-2 border-[hsl(var(--maz-warning))]/50 bg-[hsl(var(--maz-warning))]/15">
                <p class="text-base font-medium leading-relaxed text-[hsl(var(--maz-foreground))]">{{ scenario.narrative }}</p>
            </div>

            <!-- Main Content -->
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 lg:gap-8">
                <section class="flex flex-col lg:col-span-2">
                    <h2 class="text-sm font-medium uppercase tracking-wider text-[hsl(var(--maz-muted))] mb-3">
                        Scenario
                    </h2>
                    <SimulationScenarioContainer :visibility="scenario.factors.visibility as 'Clear' | 'Fog' | 'Night' | 'Rain'">
                        <div v-for="(row, rIndex) in scenario.grid_data" :key="rIndex" class="flex">
                            <SimulationCell
                                v-for="(cellCode, cIndex) in row"
                                :key="cIndex"
                                :cell-code="cellCode"
                                :definition="getCellDefinition(String(cellCode))"
                                :entities="getEntitiesAt(rIndex, cIndex)"
                                :highlight-type="getHighlightType(rIndex, cIndex)"
                                :road-condition="scenario.factors.road_condition as 'Dry' | 'Wet' | 'Icy'"
                                :lane-direction="getLaneDirection(rIndex, cIndex)"
                            />
                        </div>
                    </SimulationScenarioContainer>
                </section>

                <!-- Right: Ranking Options + Status Badges -->
                <section class="flex flex-col">
                    <h2 class="text-sm font-medium uppercase tracking-wider text-[hsl(var(--maz-muted))] mb-3">
                        Rank Options
                    </h2>
                    <ExperimentRankingOptions 
                        v-model:options="rankingOptions"
                        @highlight="highlightedZone = $event"
                        @interaction="hasUserInteracted = $event"
                    />

                    <!-- Submit Button -->
                    <div class="mt-6">
                        <MazBtn 
                            color="primary" 
                            size="lg"
                            class="w-full"
                            :disabled="store.isLoading"
                            @click="submitResponse(false)"
                        >
                            Submit Response
                        </MazBtn>
                    </div>

                    <!-- Status Badges -->
                    <div class="flex justify-center items-center flex-wrap gap-2 mt-6">
                        <span class="px-3 py-1.5 rounded-full text-sm font-medium" :class="visibilityBadgeClass">
                            {{ visibilityIcon }} {{ scenario.factors.visibility }}
                        </span>
                        <span class="px-3 py-1.5 rounded-full text-sm font-medium" :class="roadConditionBadgeClass">
                            {{ roadConditionIcon }} {{ scenario.factors.road_condition }}
                        </span>
                        <span class="px-3 py-1.5 rounded-full text-sm font-medium" :class="speedBadgeClass">
                            {{ speedIcon }} {{ scenario.factors.speed }} Speed
                        </span>
                        <span class="px-3 py-1.5 rounded-full text-sm font-medium" :class="brakeBadgeClass">
                            {{ brakeIcon }} Brake: {{ scenario.factors.brake_status }}
                        </span>
                        <span 
                            v-if="scenario.factors.has_tailgater"
                            class="px-3 py-1.5 rounded-full text-sm font-medium bg-[hsl(var(--maz-danger))]/20 text-[hsl(var(--maz-danger))]"
                        >
                            🚗 Tailgater
                        </span>
                    </div>
                </section>
            </div>
        </template>

        <!-- Error/Empty State -->
        <div v-else class="flex flex-col items-center justify-center min-h-[50vh] text-center">
            <p class="text-[hsl(var(--maz-muted))]">Failed to load scenario</p>
            <MazBtn size="sm" class="mt-4" @click="loadScenario">Retry</MazBtn>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useCellDefinition } from '~/composables/useCellDefinition'
import { useLaneDirection } from '~/composables/useLaneDirection'
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
const { getCellDefinition } = useCellDefinition()

// Timer configuration
const timerDuration = computed(() => Number(config.public.timerDuration) || 12)
const timerRef = ref()
const startTime = ref<number>(0)

// State
const scenario = ref<ScenarioResponse | null>(null)
const rankingOptions = ref<{ key: string; label: string; zone: string }[]>([])
const highlightedZone = ref<string | null>(null)
const showQuitDialog = ref(false)
const hasUserInteracted = ref(false)
const loadingState = ref<'loading' | 'submitting'>('loading')
const maxTimeoutId = ref<ReturnType<typeof setTimeout> | null>(null)
const computedInitialTime = ref(12)

// Loading message based on state
const loadingMessage = computed(() => {
    return loadingState.value === 'submitting' 
        ? 'Submitting response...' 
        : 'Loading scenario...'
})

// Highlight logic from composable
const { getHighlightType, getEntitiesAt } = useScenarioHighlight(scenario, highlightedZone)

// Lane direction from composable
const laneConfig = computed(() => scenario.value?.lane_config || { W: [], E: [], N: [], S: [] })
const { getLaneDirection } = useLaneDirection(laneConfig)

// Load scenario
async function loadScenario() {
    loadingState.value = 'loading'
    
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
        rankingOptions.value = [
            { key: 'maintain', label: data.dilemma_options.maintain, zone: 'zone_a' },
            { key: 'swerve_left', label: data.dilemma_options.swerve_left, zone: 'zone_b' },
            { key: 'swerve_right', label: data.dilemma_options.swerve_right, zone: 'zone_c' },
        ]
        
        // Set backup timeout and handle immediate expiry
        maxTimeoutId.value = setTimeout(() => {
            handleTimeUp()
        }, 15000)
        
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

// Badge helpers
const speedIcons: Record<string, string> = { Low: '🐢', Medium: '🚗', High: '🏎️' }
const brakeIcons: Record<string, string> = { Active: '✅', Fade: '⚠️', Failed: '🚨' }
const visibilityIcons: Record<string, string> = { Clear: '☀️', Fog: '🌫️', Night: '🌙', Rain: '🌧️' }
const roadConditionIcons: Record<string, string> = { Dry: '🛣️', Wet: '💧', Icy: '🧊' }

const speedIcon = computed(() => speedIcons[scenario.value?.factors.speed || ''] || '🚗')
const brakeIcon = computed(() => brakeIcons[scenario.value?.factors.brake_status || ''] || '🛞')
const visibilityIcon = computed(() => visibilityIcons[scenario.value?.factors.visibility || ''] || '☀️')
const roadConditionIcon = computed(() => roadConditionIcons[scenario.value?.factors.road_condition || ''] || '🛣️')

const badgeClasses: Record<string, string> = {
    // Speed
    Low: 'bg-[hsl(var(--maz-success))]/20 text-[hsl(var(--maz-success))]',
    Medium: 'bg-[hsl(var(--maz-warning))]/20 text-[hsl(var(--maz-warning))]',
    High: 'bg-[hsl(var(--maz-danger))]/20 text-[hsl(var(--maz-danger))]',
    // Brake
    Active: 'bg-[hsl(var(--maz-success))]/20 text-[hsl(var(--maz-success))]',
    Fade: 'bg-[hsl(var(--maz-warning))]/20 text-[hsl(var(--maz-warning))]',
    Failed: 'bg-[hsl(var(--maz-danger))]/20 text-[hsl(var(--maz-danger))]',
    // Visibility
    Clear: 'bg-[hsl(var(--maz-success))]/20 text-[hsl(var(--maz-success))]',
    Fog: 'bg-[hsl(var(--maz-warning))]/20 text-[hsl(var(--maz-warning))]',
    Night: 'bg-[hsl(var(--maz-info))]/20 text-[hsl(var(--maz-info))]',
    Rain: 'bg-[hsl(var(--maz-info))]/20 text-[hsl(var(--maz-info))]',
    // Road condition
    Dry: 'bg-[hsl(var(--maz-success))]/20 text-[hsl(var(--maz-success))]',
    Wet: 'bg-[hsl(var(--maz-warning))]/20 text-[hsl(var(--maz-warning))]',
    Icy: 'bg-[hsl(var(--maz-danger))]/20 text-[hsl(var(--maz-danger))]',
}
const defaultBadge = 'bg-[hsl(var(--maz-muted))]/20 text-[hsl(var(--maz-muted))]'

const speedBadgeClass = computed(() => badgeClasses[scenario.value?.factors.speed || ''] || defaultBadge)
const brakeBadgeClass = computed(() => badgeClasses[scenario.value?.factors.brake_status || ''] || defaultBadge)
const visibilityBadgeClass = computed(() => badgeClasses[scenario.value?.factors.visibility || ''] || defaultBadge)
const roadConditionBadgeClass = computed(() => badgeClasses[scenario.value?.factors.road_condition || ''] || defaultBadge)
</script>
