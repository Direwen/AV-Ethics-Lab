<template>
    <div class="container mx-auto px-4 py-6">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
            <div>
                <h1 class="text-xl font-semibold">{{ headerTitle }}</h1>
                <p class="text-sm text-[hsl(var(--maz-muted))]">{{ scenario.template_name }}</p>
            </div>
            <div class="flex items-center gap-4">
                <ExperimentTimer 
                    v-if="showTimer && scenario"
                    :key="scenario.id"
                    :duration="timerDuration" 
                    :initial-time="initialTime"
                    :auto-start="autoStartTimer"
                    @complete="$emit('timeUp')"
                />
                <slot name="header-actions" />
            </div>
        </div>

        <!-- Main Content -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 lg:gap-8">
            <section class="flex flex-col lg:col-span-2">
                <h2 class="text-sm font-medium uppercase tracking-wider text-[hsl(var(--maz-muted))] mb-3">
                    Scenario
                </h2>
                
                <!-- Current Conditions Dashboard -->
                <div class="mb-4 p-4 rounded-xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-card))]">
                    <div class="flex flex-wrap gap-2 justify-center items-center">
                        <!-- Speed Warning -->
                        <div class="flex items-center gap-2 px-3 py-2 rounded-lg" :class="speedWarningClass">
                            <MazBolt class="w-4 h-4" />
                            <span class="text-sm font-medium">{{ scenario.factors.speed }} Speed</span>
                        </div>
                        
                        <!-- Brake Status Warning -->
                        <div class="flex items-center gap-2 px-3 py-2 rounded-lg" :class="brakeWarningClass">
                            <MazExclamationTriangle v-if="scenario.factors.brake_status === 'Failed'" class="w-4 h-4" />
                            <MazExclamationCircle v-else-if="scenario.factors.brake_status === 'Fade'" class="w-4 h-4" />
                            <MazCheckCircle v-else class="w-4 h-4" />
                            <span class="text-sm font-medium">Brake: {{ scenario.factors.brake_status }}</span>
                        </div>
                        
                        <!-- Road Condition Warning -->
                        <div class="flex items-center gap-2 px-3 py-2 rounded-lg" :class="roadWarningClass">
                            <MazCubeTransparent v-if="scenario.factors.road_condition === 'Icy'" class="w-4 h-4" />
                            <MazCube v-else-if="scenario.factors.road_condition === 'Wet'" class="w-4 h-4" />
                            <MazSun v-else class="w-4 h-4" />
                            <span class="text-sm font-medium">{{ scenario.factors.road_condition }} Road</span>
                        </div>
                        
                        <!-- Visibility Warning -->
                        <div class="flex items-center gap-2 px-3 py-2 rounded-lg" :class="visibilityWarningClass">
                            <MazEye v-if="scenario.factors.visibility === 'Clear'" class="w-4 h-4" />
                            <MazEyeSlash v-else-if="scenario.factors.visibility === 'Fog'" class="w-4 h-4" />
                            <MazMoon v-else-if="scenario.factors.visibility === 'Night'" class="w-4 h-4" />
                            <MazCube v-else class="w-4 h-4" />
                            <span class="text-sm font-medium">Visibility: {{ scenario.factors.visibility }}</span>
                        </div>
                        
                        <!-- Tailgater Warning -->
                        <div 
                            v-if="scenario.factors.has_tailgater"
                            class="flex items-center gap-2 px-3 py-2 rounded-lg bg-[hsl(var(--maz-danger))] text-white"
                        >
                            <MazExclamationTriangle class="w-4 h-4" />
                            <span class="text-sm font-medium">Tailgater</span>
                        </div>
                    </div>
                </div>
                
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

            <section class="mt-6 lg:mt-0 flex flex-col">
                <!-- Situation Brief -->
                <h2 class="text-sm font-medium uppercase tracking-wider text-[hsl(var(--maz-muted))] mb-3">
                    Narrative
                </h2>
                <p class="text-base tracking-wider leading-relaxed text-[hsl(var(--maz-muted-foreground))]">{{ scenario.narrative }}</p>

                <div class="border border-[hsl(var(--maz-muted))] w-full my-6"></div>

                <!-- Layer 3: Ranking Controls -->
                <h2 class="text-sm font-medium uppercase tracking-wider text-[hsl(var(--maz-muted))] mb-3">
                    {{ rankingTitle }}
                </h2>
                <ExperimentRankingOptions 
                    v-model:options="localRankingOptions"
                    @highlight="$emit('highlight', $event)"
                    @interaction="$emit('interaction', $event)"
                />

                <!-- Layer 4: Submit Action -->
                <div class="mt-6">
                    <MazBtn 
                        color="primary" 
                        size="lg"
                        class="w-full"
                        :disabled="isLoading"
                        @click="$emit('submit')"
                    >
                        {{ submitButtonText }}
                    </MazBtn>
                </div>
            </section>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useCellDefinition } from '~/composables/useCellDefinition'
import { useLaneDirection } from '~/composables/useLaneDirection'
import type { ScenarioResponse } from '~/types/response.types'
import { 
    MazBolt, 
    MazExclamationTriangle, 
    MazExclamationCircle, 
    MazCheckCircle,
    MazCubeTransparent,
    MazCube,
    MazSun,
    MazEye,
    MazEyeSlash,
    MazMoon
} from '@maz-ui/icons'

export interface RankingOption {
    key: string
    label: string
    zone: string
}

const props = withDefaults(defineProps<{
    scenario: ScenarioResponse
    rankingOptions: RankingOption[]
    highlightedZone?: string | null
    isLoading?: boolean
    showTimer?: boolean
    timerDuration?: number
    initialTime?: number
    autoStartTimer?: boolean
    headerTitle?: string
    rankingTitle?: string
    submitButtonText?: string
}>(), {
    highlightedZone: null,
    isLoading: false,
    showTimer: true,
    timerDuration: 20,
    initialTime: 20,
    autoStartTimer: true,
    headerTitle: 'Scenario',
    rankingTitle: 'Which action are you likely to do? Rank these?',
    submitButtonText: 'Submit Response'
})

const emit = defineEmits<{
    'update:rankingOptions': [options: RankingOption[]]
    'highlight': [zone: string | null]
    'interaction': [hasInteracted: boolean]
    'submit': []
    'timeUp': []
}>()

const { getCellDefinition } = useCellDefinition()

// Local copy of ranking options for v-model
const localRankingOptions = ref([...props.rankingOptions])

// Watch for external changes - only sync if arrays are actually different
watch(() => props.rankingOptions, (newVal) => {
    const newKeys = newVal.map(o => o.key).join(',')
    const localKeys = localRankingOptions.value.map(o => o.key).join(',')
    if (newKeys !== localKeys) {
        localRankingOptions.value = [...newVal]
    }
}, { deep: true })

// Emit changes when ranking options change - only if different from props
watch(localRankingOptions, (newVal) => {
    const newKeys = newVal.map(o => o.key).join(',')
    const propKeys = props.rankingOptions.map(o => o.key).join(',')
    if (newKeys !== propKeys) {
        emit('update:rankingOptions', [...newVal])
    }
}, { deep: true })

// Highlight logic from composable
const { getHighlightType, getEntitiesAt } = useScenarioHighlight(toRef(props, 'scenario'), toRef(props, 'highlightedZone'))

// Lane direction from composable
const laneConfig = computed(() => props.scenario?.lane_config || { W: [], E: [], N: [], S: [] })
const { getLaneDirection } = useLaneDirection(laneConfig)

// Warning Dashboard Classes (Automotive-style)
const speedWarningClass = computed(() => {
    const speed = props.scenario?.factors.speed
    switch (speed) {
        case 'High': return 'bg-[hsl(var(--maz-danger))] text-white'
        case 'Medium': return 'bg-[hsl(var(--maz-warning))] text-[hsl(var(--maz-warning-foreground))]'
        case 'Low': return 'bg-[hsl(var(--maz-success))]/20 text-[hsl(var(--maz-success))]'
        default: return 'bg-[hsl(var(--maz-muted))]/20 text-[hsl(var(--maz-muted))]'
    }
})

const brakeWarningClass = computed(() => {
    const brake = props.scenario?.factors.brake_status
    switch (brake) {
        case 'Failed': return 'bg-[hsl(var(--maz-danger))] text-white'
        case 'Fade': return 'bg-[hsl(var(--maz-warning))] text-[hsl(var(--maz-warning-foreground))]'
        case 'Active': return 'bg-[hsl(var(--maz-success))]/20 text-[hsl(var(--maz-success))]'
        default: return 'bg-[hsl(var(--maz-muted))]/20 text-[hsl(var(--maz-muted))]'
    }
})

const roadWarningClass = computed(() => {
    const road = props.scenario?.factors.road_condition
    switch (road) {
        case 'Icy': return 'bg-[hsl(var(--maz-danger))] text-white'
        case 'Wet': return 'bg-[hsl(var(--maz-warning))] text-[hsl(var(--maz-warning-foreground))]'
        case 'Dry': return 'bg-[hsl(var(--maz-success))]/20 text-[hsl(var(--maz-success))]'
        default: return 'bg-[hsl(var(--maz-muted))]/20 text-[hsl(var(--maz-muted))]'
    }
})

const visibilityWarningClass = computed(() => {
    const visibility = props.scenario?.factors.visibility
    switch (visibility) {
        case 'Fog': return 'bg-[hsl(var(--maz-danger))] text-white'
        case 'Night': 
        case 'Rain': return 'bg-[hsl(var(--maz-warning))] text-[hsl(var(--maz-warning-foreground))]'
        case 'Clear': return 'bg-[hsl(var(--maz-success))]/20 text-[hsl(var(--maz-success))]'
        default: return 'bg-[hsl(var(--maz-muted))]/20 text-[hsl(var(--maz-muted))]'
    }
})
</script>