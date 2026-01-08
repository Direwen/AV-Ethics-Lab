<template>
    <div class="container mx-auto px-4 py-12">
        <div class="max-w-2xl mx-auto text-center">
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
                Your responses will contribute valuable insights to the development of ethical AI systems.
            </p>

            <!-- Completion Stats -->
            <div class="bg-[hsl(var(--maz-background-accent))]/50 rounded-xl p-6 mb-8">
                <h2 class="text-lg font-semibold mb-4 text-[hsl(var(--maz-foreground))]">
                    Experiment Completed
                </h2>
                <div class="flex justify-center items-center gap-2 text-[hsl(var(--maz-success))]">
                    <MazIcon name="check-circle" class="w-5 h-5" />
                    <span class="font-medium">All scenarios completed successfully</span>
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
                    class="custom-background-gradient "
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
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { MazCheck } from '@maz-ui/icons'

definePageMeta({
    layout: 'default'
})

const router = useRouter()
const toast = useMazToast()

// Optional: Add some confetti or celebration animation
onMounted(() => {
    // You could add a celebration animation here
    toast.success('Experiment completed successfully!', {
        timeout: 5000
    })
})
</script>