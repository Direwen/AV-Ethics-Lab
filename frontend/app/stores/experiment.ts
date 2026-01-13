import { defineStore } from "pinia";
import type { Answer } from "~/types/answer.types";
import type { Demographic } from "~/types/demographic.types";
import type {
    SubmitResponsePayload
} from "~/types/request.types"
import type { 
    ApiResponse, 
    CreateSessionResponse, 
    ScenarioResponse,
    ResponseSubmissionResult,
    FeedbackResponse,
    DashboardStats
} from "~/types/response.types";
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
            const response = await $api<ApiResponse<CreateSessionResponse>>('/api/v1/sessions', {
                method: 'POST',
                body: {
                    ...demographics,
                    fingerprint: await getFingerprint(),
                    self_reported_new: selfReportedNew
                }
            })

            if (!response.success) {
                throw new Error(response.message)
            }

            token.value = response.data?.token
            await nextTick()

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

    
    async function getScenario() {
        isLoading.value = true
        try {
            const response = await $api<ApiResponse<ScenarioResponse>>('/api/v1/scenarios/next', {method: 'GET'})

            if (!response.success) {
                throw new Error(response.message)
            }
            return response.data

        } catch(e: any) {
            toast.error("Failed to get the scenario")
        } finally {
            isLoading.value = false
        }
    }

    async function submitResponse(scenarioId: string, payload: SubmitResponsePayload) {
        isLoading.value = true
        try {
            const response = await $api<ApiResponse<ResponseSubmissionResult>>(`/api/v1/scenarios/${scenarioId}/responses`, {
                method: 'POST',
                body: payload
            })

            if (!response.success) {
                throw new Error(response.message)
            }

            toast.success("Response submitted successfully")
            return response.data

        } catch(e: any) {
            toast.error("Failed to submit response")
            throw e
        } finally {
            isLoading.value = false
        }
    }

    async function getFeedback() {
        isLoading.value = true
        try {
            const response = await $api<ApiResponse<FeedbackResponse>>('/api/v1/sessions/feedback', {method: 'GET'})
            if (!response.success) {
                throw new Error(response.message)
            }
            return response.data
        } catch (e: any) {
            return {
                archetype: 'The Thoughtful Participant',
                summary: 'Thank you for completing the experiment. We were unable to generate your personalized feedback at this time, but your responses have been recorded and will contribute to our research.',
                key_trait: 'Valued Contributor'
            }
        } finally {
            isLoading.value = false
        }
    }

    async function getDashboardData() {
        isLoading.value = true
        try {
            const response = await $api<ApiResponse<DashboardStats>>('/api/v1/dashboard', {method: 'GET'})
            if (!response.success) {
                throw new Error(response.message)
            }
            toast.success("Dashboard stats fetched successfully")
            return response.data
        } catch (e: any) {
            toast.error("Failed to get dashboard stats")
            return
        } finally {
            isLoading.value = false
        }
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
        getFingerprint,
        getScenario,
        submitResponse,
        getFeedback,
        getDashboardData
    }
})