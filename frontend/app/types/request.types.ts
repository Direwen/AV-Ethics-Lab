import type { Demographic } from "./demographic.types";

export interface CreateSessionRequest {
    demographics: Demographic,
    fingerprint: string,
    selfReportedNew: boolean
}
export interface SubmitResponsePayload {
    ranking_order: string[]
    response_time_ms: number
    is_timeout: boolean
    has_interacted: boolean
}