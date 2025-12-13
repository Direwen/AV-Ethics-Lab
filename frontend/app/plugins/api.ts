import { useNuxtApp } from '#app'
import { useErrorParser } from '~/composables/useErrorParser'

export default defineNuxtPlugin((nuxtApp) => {
    const config = useRuntimeConfig()
    const { parseError } = useErrorParser()

    const api = $fetch.create({
        baseURL: config.public.apiBase as string,

        onRequest({ options }) {
            const token = useCookie('token')
            if (token.value) {
                // Ensure headers is a Headers object
                const headers = new Headers(options.headers)
                headers.set('Authorization', `Bearer ${token.value}`)
                options.headers = headers
            }
        },
        
        onResponseError({ response }) {
            if (response.status === 401) {
                const token = useCookie('token')
                token.value = null
                navigateTo('/auth/login')
            }
        }
    })

    return {
        provide: {
            api
        }
    }
})

declare module '#app' {
    interface NuxtApp {
        $api: typeof $fetch
    }
}

declare module 'vue' {
    interface ComponentCustomProperties {
        $api: typeof $fetch
    }
}