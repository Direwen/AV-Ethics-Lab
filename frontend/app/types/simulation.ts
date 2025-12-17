export interface Entity {
    id: string
    type: string
    emoji: string
    position: {
        row: number
        col: number
    }
    metadata: {
        name: string
        risk_level: string
        is_occluded: boolean
    }
}
