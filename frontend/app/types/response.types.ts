import type { Entity } from './simulation'

export interface ApiResponse<T = unknown> {
    success: boolean
    message: string
    data?: T
    error?: string
}

export interface CreateSessionResponse {
    token: string
}

export interface DilemmaOptions {
    maintain: string
    swerve_left: string
    swerve_right: string
}

export interface ScenarioFactors {
    visibility: 'Clear' | 'Fog' | 'Night' | 'Rain'
    road_condition: 'Dry' | 'Wet' | 'Icy'
    location: string
    brake_status: 'Active' | 'Fade' | 'Failed'
    speed: 'Low' | 'Medium' | 'High'
    has_tailgater: boolean
    primary_entity: string
    primary_behavior: string
    background_entities: string[]
}

export interface ZoneCoordinate {
    row: number
    col: number
    surface: 'drivable' | 'walkable'
    orientation: string
}

export interface TridentZone {
    coordinates: ZoneCoordinate[]
}

export interface TridentZones {
    zone_a: TridentZone
    zone_b: TridentZone
    zone_c: TridentZone
}

export interface LaneConfig {
    E: number[][]
    N: number[][]
    S: number[][]
    W: number[][]
}

export interface ScenarioResponse {
    id: string
    narrative: string
    dilemma_options: DilemmaOptions
    entities: Entity[]
    factors: ScenarioFactors
    width: number
    height: number
    grid_data: number[][]
    lane_config: LaneConfig
    trident_zones: TridentZones
    template_name: string
    current_step: number
    total_steps: number
}

export interface ResponseSubmissionResult {
    id: string
    scenario_id: string
    ranking_order: string[]
    response_time_ms: number
    is_timeout: boolean
    has_interacted: boolean
    created_at: string
}

export interface FeedbackResponse {
    archetype: string,
    summary: string,
    key_trait: string
}

// Dashboard Types
export interface DashboardStats {
    completed_sessions: number
    countries_represented: number
    least_harmful_outcome: OutcomeDistribution
    self_preservation_effect: TailgaterEffect
    entity_compliance_effect: ComplianceEffect
    decision_time_distribution: TimeDistributionPoint[]
    archetype_distribution: ArchetypeCount[]
}

export interface OutcomeDistribution {
    maintain: number
    swerve_left: number
    swerve_right: number
    total: number
    maintain_pct: number
    swerve_left_pct: number
    swerve_right_pct: number
}

export interface TailgaterEffect {
    with_tailgater: EffectMetric
    without_tailgater: EffectMetric
}

export interface EffectMetric {
    maintain_count: number
    total_count: number
    percentage: number
}

export interface ComplianceEffect {
    compliant: EffectMetric
    violation: EffectMetric
}

export interface TimeDistributionPoint {
    seconds: number
    count: number
}

export interface ArchetypeCount {
    archetype: string
    count: number
}
