import { defineStore } from "pinia";
import type { Answer } from "~/types/answer.types";
import type { Demographic } from "~/types/demographic.types";
import type { CreateSessionRequest } from "~/types/request.types";
import type { ApiResponse, CreateSessionResponse } from "~/types/response.types";
import type { Scenario } from "~/types/scenario.types";

export const useExperimentStore = defineStore('experiment', () => {
    const { $api, $fingerprint } = useNuxtApp()
    const toast = useMazToast()

    const token = useCookie('session_token', { maxAge: 60 * 60 * 4 }) // 4 Hours    
    const fingerprint = useCookie('fingerprint')
    const isLoading = ref(false)

    const scenarios = ref<Scenario[]>([])
    const currentScenarioIndex = ref<number>(0)
    const currentScenario = ref<Scenario | null>(null)

    const answers = ref<Answer[]>([])
    const demographic = ref<Demographic>({} as Demographic)

    async function createSession(demographics: Demographic, selfReportedNew: boolean) {
        isLoading.value = true
        try {
            const response = await $api<ApiResponse<CreateSessionResponse>>('', {
                method: 'POST',
                body: {
                    demographics,
                    fingerprint: await getFingerprint(),
                    selfReportedNew
                } as CreateSessionRequest
            })

            if (response.success) {
                token.value = response.data?.token
                await nextTick()
            }
        } catch(e: any) {
            toast.error("Failed to Create Session")
        } finally {
            isLoading.value = false
        }
    }

    async function init() {
        if (!token.value) return
        isLoading.value = false
        try {
            
        } catch (e) {
            
        } finally {
            isLoading.value = false
        }
    }

    async function getFingerprint() {
        if (fingerprint.value) return fingerprint.value
        const result = await $fingerprint.get()
        fingerprint.value = result.visitorId
        return result.visitorId
    }

    return {
        // State
        token,
        isLoading,
        scenarios,
        currentScenarioIndex,
        currentScenario,
        answers,
        demographic,
        
        // Actions
        createSession,
        init,
        getFingerprint
    }
})