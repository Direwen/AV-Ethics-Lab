<template>
    <div>
        <ExperimentScenarioViewer
            :scenario="staticScenario"
            :ranking-options="rankingOptions"
            :highlighted-zone="highlightedZone"
            :is-loading="false"
            :show-timer="false"
            header-title="Calibration Test"
            ranking-title="Which action are you likely to do?"
            submit-button-text="Continue"
            @highlight="highlightedZone = $event"
            @interaction="hasUserInteracted = $event"
            @submit="handleSubmit"
        />

        <!-- Tutorial Dialog -->
        <MazDialog v-model="showTutorial" title="How to Participate">
            <div class="space-y-4">
                <div>
                    <h3 class="font-semibold mb-2">Welcome to the Experiment!</h3>
                    <p class="text-sm text-[hsl(var(--maz-muted-foreground))]">
                        This is a practice scenario to help you understand how the experiment works.
                    </p>
                </div>
                
                <div>
                    <h4 class="font-medium mb-1">1. Review the Scenario</h4>
                    <p class="text-sm text-[hsl(var(--maz-muted-foreground))]">
                        Look at the driving conditions (speed, road, visibility, brakes) and read the situation description.
                    </p>
                </div>
                
                <div>
                    <h4 class="font-medium mb-1">2. Study the Visual</h4>
                    <p class="text-sm text-[hsl(var(--maz-muted-foreground))]">
                        The grid shows the road layout with vehicles, pedestrians, and other elements.
                    </p>
                </div>
                
                <div>
                    <h4 class="font-medium mb-1">3. Rank Your Choices</h4>
                    <p class="text-sm text-[hsl(var(--maz-muted-foreground))]">
                        Drag and drop the action options to rank them from most likely (top) to least likely (bottom) that you would choose.
                    </p>
                </div>
                
                <div>
                    <h4 class="font-medium mb-1">4. Hover to Preview</h4>
                    <p class="text-sm text-[hsl(var(--maz-muted-foreground))]">
                        Hover over ranking options to see highlighted zones on the grid showing where each action would take you.
                    </p>
                </div>
                
                <div class="p-3 rounded-lg bg-[hsl(var(--maz-info))]/10 border border-[hsl(var(--maz-info))]/20">
                    <p class="text-sm font-medium text-[hsl(var(--maz-info))]">
                        💡 Take your time to explore this practice scenario. In the real experiment, you'll have a time limit!
                    </p>
                </div>
            </div>
            
            <template #footer="{ close }">
                <div class="flex gap-3 justify-end">
                    <MazBtn color="primary" @click="close">Got it, let's practice!</MazBtn>
                </div>
            </template>
        </MazDialog>
    </div>
</template>

<script setup lang="ts">
import type { ScenarioResponse } from '~/types/response.types'
import { useExperimentStore } from '~/stores/experiment'

definePageMeta({
    layout: 'optional-layout',
    // middleware: ['guide']
})

// Static scenario data
const staticScenario: ScenarioResponse = {
    id: "a54a6291-09b0-4a68-b685-88ab2ae4eacf",
    narrative: "An icy intersection dilemma: brake for a jaywalking cat or swerve into traffic/pedestrians.",
    dilemma_options: {
        maintain: "Maintain Course to Avoid Cat, Risk Braking on Ice",
        swerve_left: "Swerve Left to Avoid Cat, Risk Hitting Sports Car or Jogger",
        swerve_right: "Swerve Right to Avoid Cat, Risk Hitting Car or Pedestrian"
    },
    entities: [
        {
            id: "ent_vehicle_av_ego",
            type: "vehicle_av",
            emoji: "🚕",
            row: 4,
            col: 13,
            metadata: {
                is_star: false,
                is_ego: true,
                is_violation: false,
                action: "",
                orientation: "W"
            }
        },
        {
            id: "ent_animal_cat_0",
            type: "animal_cat",
            emoji: "🐈",
            row: 4,
            col: 11,
            metadata: {
                is_star: true,
                is_ego: false,
                is_violation: true,
                action: "Jaywalking across road",
                orientation: "W"
            }
        },
        {
            id: "ent_vehicle_sports_car_1",
            type: "vehicle_sports_car",
            emoji: "🏎️",
            row: 5,
            col: 11,
            metadata: {
                is_star: false,
                is_ego: false,
                is_violation: false,
                action: "Approaching intersection",
                orientation: "N"
            }
        },
        {
            id: "ent_ped_jogger_2",
            type: "ped_jogger",
            emoji: "🏃‍♀️",
            row: 5,
            col: 9,
            metadata: {
                is_star: false,
                is_ego: false,
                is_violation: true,
                action: "Jaywalking across road",
                orientation: "S"
            }
        },
        {
            id: "ent_vehicle_car_3",
            type: "vehicle_car",
            emoji: "🚗",
            row: 3,
            col: 11,
            metadata: {
                is_star: false,
                is_ego: false,
                is_violation: false,
                action: "Stopped at intersection",
                orientation: "N"
            }
        },
        {
            id: "ent_ped_adult_4",
            type: "ped_adult",
            emoji: "🧍",
            row: 2,
            col: 10,
            metadata: {
                is_star: false,
                is_ego: false,
                is_violation: true,
                action: "Jaywalking into street",
                orientation: ""
            }
        }
    ],
    factors: {
        visibility: "Clear",
        road_condition: "Icy",
        location: "CN",
        brake_status: "Active",
        speed: "High",
        has_tailgater: false,
        primary_entity: "animal_cat",
        primary_behavior: "Violation",
        background_entities: ["ped_jogger", "ped_business", "obstacle_barrier", "vehicle_car", "ped_business", "ped_jogger", "vehicle_sports_car", "ped_business", "ped_adult", "vehicle_car"]
    },
    width: 20,
    height: 11,
    grid_data: [
        [0,0,0,0,0,0,0,0,19,10,17,10,20,0,0,0,0,0,0,0],
        [0,0,0,0,0,0,0,0,19,10,17,10,20,0,0,0,0,0,0,0],
        [1,1,1,1,1,1,1,1,19,15,15,15,20,1,1,1,1,1,1,1],
        [3,3,3,3,3,3,3,3,5,11,17,11,7,3,3,3,3,3,3,3],
        [9,9,9,9,9,9,9,16,11,11,11,11,11,16,9,9,9,9,9,9],
        [18,18,18,18,18,18,18,16,18,11,11,11,18,16,18,18,18,18,18,18],
        [9,9,9,9,9,9,9,16,11,11,11,11,11,16,9,9,9,9,9,9],
        [4,4,4,4,4,4,4,4,6,11,17,11,8,4,4,4,4,4,4,4],
        [2,2,2,2,2,2,2,2,19,15,15,15,20,2,2,2,2,2,2,2],
        [0,0,0,0,0,0,0,0,19,10,17,10,20,0,0,0,0,0,0,0],
        [0,0,0,0,0,0,0,0,19,10,17,10,20,0,0,0,0,0,0,0]
    ],
    lane_config: {
        E: [[6,0],[6,1],[6,2],[6,3],[6,4],[6,5],[6,6],[6,7],[6,8],[6,9],[6,10],[6,11],[6,12],[6,13],[6,14],[6,15],[6,16],[6,17],[6,18],[6,19]],
        N: [[0,11],[1,11],[2,11],[3,11],[4,11],[5,11],[6,11],[7,11],[8,11],[9,11],[10,11]],
        S: [[0,9],[1,9],[2,9],[3,9],[4,9],[5,9],[6,9],[7,9],[8,9],[9,9],[10,9]],
        W: [[4,0],[4,1],[4,2],[4,3],[4,4],[4,5],[4,6],[4,7],[4,8],[4,9],[4,10],[4,11],[4,12],[4,13],[4,14],[4,15],[4,16],[4,17],[4,18],[4,19]]
    },
    trident_zones: {
        zone_a: {
            coordinates: [
                { row: 4, col: 11, surface: "drivable", orientation: "W" },
                { row: 4, col: 10, surface: "drivable", orientation: "W" },
                { row: 4, col: 9, surface: "drivable", orientation: "W" }
            ]
        },
        zone_b: {
            coordinates: [
                { row: 5, col: 11, surface: "drivable", orientation: "N" },
                { row: 5, col: 10, surface: "drivable", orientation: "" },
                { row: 5, col: 9, surface: "drivable", orientation: "S" }
            ]
        },
        zone_c: {
            coordinates: [
                { row: 3, col: 11, surface: "drivable", orientation: "N" },
                { row: 2, col: 10, surface: "drivable", orientation: "" },
                { row: 3, col: 9, surface: "drivable", orientation: "S" }
            ]
        }
    },
    template_name: "4-Way Urban Intersection",
    current_step: 1,
    total_steps: 2
}

// State
const showTutorial = ref(true) // Auto-open tutorial dialog
const highlightedZone = ref<string | null>(null)
const hasUserInteracted = ref(false)

// Initialize ranking options with randomization
const rankingOptions = ref((() => {
    const options = [
        { key: 'maintain', label: staticScenario.dilemma_options.maintain, zone: 'zone_a' },
        { key: 'swerve_left', label: staticScenario.dilemma_options.swerve_left, zone: 'zone_b' },
        { key: 'swerve_right', label: staticScenario.dilemma_options.swerve_right, zone: 'zone_c' },
    ]
    
    // Randomize to prevent passive agreement bias
    for (let i = options.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [options[i], options[j]] = [options[j]!, options[i]!]
    }
    
    return options
})())

function handleSubmit() {
    console.log('Calibration response:', {
        ranking_order: rankingOptions.value.map(opt => opt.key),
        has_interacted: hasUserInteracted.value
    })
    
    // Mark guide as completed and redirect to experiment
    const store = useExperimentStore()
    store.completeGuide()
    
    // Navigate to the actual experiment
    navigateTo('/experiment')
}
</script>