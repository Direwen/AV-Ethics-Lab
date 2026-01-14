export default defineNuxtRouteMiddleware((to) => {
    const store = useExperimentStore()
    
    // If user doesn't have a session token, redirect to consent
    if (!store.token) {
        return navigateTo('/experiment/consent')
    }
    
    // If user can't access guide (returning user or already completed), redirect to experiment
    if (!store.canAccessGuide) {
        return navigateTo('/experiment')
    }
})