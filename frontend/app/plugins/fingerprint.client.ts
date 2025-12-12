import FingerprintJS from '@fingerprintjs/fingerprintjs'
import type { Agent } from '@fingerprintjs/fingerprintjs'

export default defineNuxtPlugin(async () => {
      // Load the agent at application startup.
    const fp = await FingerprintJS.load()

    // Return the agent instance for injection
    return {
        provide: {
        fingerprint: fp
        }
    }
})

declare module '#app' {
    interface NuxtApp {
        $fingerprint: Agent
    }
}

declare module 'vue' {
    interface ComponentCustomProperties {
        $fingerprint: Agent
    }
}

export {}

