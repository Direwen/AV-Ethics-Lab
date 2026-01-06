export default defineNuxtRouteMiddleware((to) => {
    const token = useCookie('session_token')
    const protectedRoutes = ['/experiment', '/feedback']
    const guestRoutes = ['/', '/experiment/consent']

    // If trying to access protected route WITHOUT token, redirect to home
    if (protectedRoutes.includes(to.path) && !token.value) {
        return navigateTo('/')
    }

    // If trying to access guest route WITH token, redirect to experiment (FOR NOW)
    if (guestRoutes.includes(to.path) && token.value) {
        return navigateTo('/experiment')
    }
})