export default defineNuxtRouteMiddleware((to) => {
    const token = useCookie('session_token')
    const protectedRoutes = ['/experiment', '/feedback']
    const guestRoutes = ['/', '/experiment/consent']

    const toast = useMazToast()


    // If trying to access protected route WITHOUT token, redirect to home
    if (protectedRoutes.includes(to.path) && !token.value) {
        toast.error("Please submit the consent form first")
        return navigateTo('/')
    }

    // If trying to access guest route WITH token, redirect to experiment (FOR NOW)
    if (guestRoutes.includes(to.path) && token.value) {
        toast.error("Please Quit the ongoing experiment first")
        return navigateTo('/experiment')
    }
})