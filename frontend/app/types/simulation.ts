export interface Position {
    row: number
    col: number
}

export interface EntityMetadata {
    is_star: boolean
    is_ego: boolean
    is_violation: boolean
    action: string
    orientation: string
}

export interface Entity {
    id: string
    type: string
    emoji: string
    row: number
    col: number
    metadata: EntityMetadata
}

export interface ScenarioData {
    scenario_id: string
    template_meta: {
        id: string
        name: string
        dimensions: { rows: number; cols: number }
    }
    grid: number[][]
    entities: Entity[]
}
export interface CellDefinition {
    name?: string
    class: string
    isInteractive: boolean
}