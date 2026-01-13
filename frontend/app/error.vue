<template>
    <div class="min-h-screen flex items-center justify-center bg-radial from-[hsl(var(--maz-background-accent))] to-[hsl(var(--maz-background))]">
        <div class="text-center px-4">
        <!-- 404 Visual -->
        <div class="mb-8">
            <h1 class="text-[8rem] md:text-[12rem] font-bold leading-none custom-background-gradient bg-clip-text text-transparent">
            {{ error?.statusCode || 404 }}
            </h1>
        </div>

        <!-- Message -->
        <h2 class="text-2xl md:text-3xl font-semibold mb-4 text-[hsl(var(--maz-foreground))]">
            {{ errorTitle }}
        </h2>
        <p class="text-[hsl(var(--maz-foreground))]/60 max-w-md mx-auto mb-8">
            {{ errorMessage }}
        </p>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <MazBtn 
            color="primary"
            size="lg"
            class="custom-background-gradient"
            @click="handleError"
            >
            <MazHome class="w-5 h-5 mr-2" />
            Back to Home
            </MazBtn>
            <MazBtn 
                color="background"
                variant="outline"
                size="lg"
                @click="goBack"
                class=""
            >
                <MazArrowLeft class="w-5 h-5 mr-2" />
                Go Back
            </MazBtn>
        </div>

        <!-- Decorative wheels -->
        <div class="mt-16 flex justify-center gap-8 opacity-20">
            <div class="wheel-container animate-spin" style="animation-duration: 3s;">
            <div class="wheel-outer custom-background-gradient" />
            <div class="wheel-inner bg-[hsl(var(--maz-background))]">
                <div class="wheel-hub custom-background-gradient opacity-50" />
            </div>
            </div>
            <div class="wheel-container animate-spin" style="animation-duration: 3s;">
            <div class="wheel-outer custom-background-gradient" />
            <div class="wheel-inner bg-[hsl(var(--maz-background))]">
                <div class="wheel-hub custom-background-gradient opacity-50" />
            </div>
            </div>
        </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { MazHome, MazArrowLeft } from '@maz-ui/icons'

const props = defineProps<{
    error: {
        statusCode?: number
        message?: string
    }
}>()

const router = useRouter()

const errorTitle = computed(() => {
    switch (props.error?.statusCode) {
        case 404: return 'Page Not Found'
        case 500: return 'Server Error'
        case 403: return 'Access Denied'
        default: return 'Something Went Wrong'
    }
})

const errorMessage = computed(() => {
    switch (props.error?.statusCode) {
        case 404: return "The page you're looking for doesn't exist or has been moved."
        case 500: return 'Our servers encountered an unexpected error. Please try again later.'
        case 403: return "You don't have permission to access this page."
        default: return props.error?.message || 'An unexpected error occurred.'
    }
})

function handleError() {
    clearError({ redirect: '/' })
}

function goBack() {
    router.back()
}
</script>

<style scoped>
.wheel-container {
    position: relative;
    width: 48px;
    height: 48px;
}

.wheel-outer {
    position: absolute;
    inset: 0;
    border-radius: 50%;
    mask: 
        linear-gradient(#fff 0 0) content-box,
        linear-gradient(#fff 0 0);
    mask-composite: exclude;
    padding: 4px;
}

.wheel-inner {
    position: absolute;
    inset: 4px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
}

.wheel-hub {
    width: 12px;
    height: 12px;
    border-radius: 50%;
}
</style>
