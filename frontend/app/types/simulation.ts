export interface Position {
    row: number
    col: number
}

export interface EntityMetadata {
    name: string
    risk_level: string 
    is_occluded: boolean
    behavior?: string
    occluded_by?: string // ID of the entity hiding this one
}

export interface Entity {
    id: string
    type: string
    emoji: string
    position: Position
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