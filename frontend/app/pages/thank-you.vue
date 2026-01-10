<template>
    <div class="container mx-auto px-4 py-12">
        <div class="max-w-2xl mx-auto text-center">
            <!-- Loading State -->
            <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[50vh]">
                <MazSpinner size="3em" />
                <p class="mt-4 text-lg text-[hsl(var(--maz-muted))]">Generating your personalized feedback...</p>
                <p class="mt-2 text-sm text-[hsl(var(--maz-muted))]/70">Analyzing your decision patterns</p>
            </div>

            <!-- Feedback Content -->
            <template v-else>
                <!-- Success Icon -->
                <div class="mb-8">
                    <div class="w-24 h-24 mx-auto bg-[hsl(var(--maz-success))]/20 rounded-full flex items-center justify-center">
                        <MazCheck class="w-12 h-12 text-[hsl(var(--maz-success))]" />
                    </div>
                </div>

                <!-- Thank You Message -->
                <h1 class="text-3xl font-bold mb-4 text-[hsl(var(--maz-foreground))]">
                    Thank You for Participating!
                </h1>
                
                <p class="text-lg text-[hsl(var(--maz-muted))] mb-8 leading-relaxed">
                    You have successfully completed the autonomous vehicle ethics experiment.
                </p>

                <!-- Personalized Feedback Card -->
                <div v-if="feedback" class="bg-[hsl(var(--maz-background-accent))]/50 rounded-xl p-6 mb-8 text-left">
                    <h2 class="text-lg font-semibold mb-4 text-[hsl(var(--maz-foreground))] text-center">
                        Your Ethical Profile
                    </h2>
                    
                    <!-- Archetype -->
                    <div class="mb-4 text-center">
                        <span class="inline-block px-4 py-2 rounded-full bg-[hsl(var(--maz-primary))]/20 text-[hsl(var(--maz-primary))] font-semibold text-lg">
                            {{ feedback.archetype }}
                        </span>
                    </div>

                    <!-- Key Trait -->
                    <div class="mb-4 text-center">
                        <span class="text-sm text-[hsl(var(--maz-muted))]">Key Trait:</span>
                        <span class="ml-2 font-medium text-[hsl(var(--maz-foreground))]">{{ feedback.key_trait }}</span>
                    </div>

                    <!-- Summary -->
                    <div class="border-t border-[hsl(var(--maz-border))] pt-4 mt-4">
                        <p class="text-[hsl(var(--maz-foreground))] leading-relaxed">
                            {{ feedback.summary }}
                        </p>
                    </div>
                </div>

                <!-- Additional Information -->
                <div class="text-sm text-[hsl(var(--maz-muted))] mb-8 space-y-2">
                    <p>Your session has been automatically closed for privacy and security.</p>
                    <p>If you have any questions about this research, please contact the research team.</p>
                </div>

                <!-- Action Buttons -->
                <div class="flex flex-col sm:flex-row gap-4 justify-center">
                    <MazBtn 
                        color="primary" 
                        class="custom-background-gradient"
                        size="lg"
                        @click="router.push('/')"
                    >
                        Return to Home
                    </MazBtn>
                </div>

                <!-- Research Information -->
                <div class="mt-12 pt-8 border-t border-[hsl(var(--maz-border))]">
                    <h3 class="text-lg font-semibold mb-4 text-[hsl(var(--maz-foreground))]">
                        About This Research
                    </h3>
                    <p class="text-sm text-[hsl(var(--maz-muted))] leading-relaxed">
                        This experiment is part of ongoing research into ethical decision-making frameworks 
                        for autonomous vehicles. Your participation helps us understand human moral reasoning 
                        in complex traffic scenarios, which will inform the development of more ethical AI systems.
                    </p>
                </div>
            </template>
        </div>
    </div>
</template>

<script setup lang="ts">
import { MazCheck } from '@maz-ui/icons'
import { useExperimentStore } from '~/stores/experiment'
import type { FeedbackResponse } from '~/types/response.types'

definePageMeta({
    layout: 'default'
})

const router = useRouter()
const store = useExperimentStore()

const isLoading = ref(true)
const feedback = ref<FeedbackResponse | null>(null)

onMounted(async () => {
    feedback.value = await store.getFeedback() || null
    isLoading.value = false
    
    // Clear session after feedback is fetched
    store.token = null
})
</script>
