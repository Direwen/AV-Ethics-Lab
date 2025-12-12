import type { Demographic } from "./demographic.types";

export interface CreateSessionRequest {
    demographics: Demographic,
    fingerprint: string,
    selfReportedNew: boolean
}