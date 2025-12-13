export default defineNuxtRouteMiddleware((to, from) => {
    const token = useCookie('token')
    const protectedRoutes = ['/experiment', '/feedback']
    const guestRoutes = ['/', '/auth/consent']

    // If trying to access protected route WITHOUT token, redirect to home
    if (protectedRoutes.includes(to.path) && !token.value) {
        return navigateTo('/')
    }

    // If trying to access guest route WITH token, redrect to experiment (FOR NOW)
    if (guestRoutes.includes(to.path) && token.value) {
        return navigateTo('/experiment')
    }
})