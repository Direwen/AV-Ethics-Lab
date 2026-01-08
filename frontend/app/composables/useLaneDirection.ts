import { computed, type ComputedRef } from 'vue'

export interface LaneConfig {
    W: number[][]
    E: number[][]
    N: number[][]
    S: number[][]
}

export function useLaneDirection(laneConfig: ComputedRef<LaneConfig>) {
    const laneDirectionMap = computed(() => {
        const map = new Map<string, string>()
        
        for (const [direction, cells] of Object.entries(laneConfig.value)) {
        for (const [row, col] of cells as number[][]) {
            map.set(`${row},${col}`, direction)
        }
        }
        return map
    })

    function getLaneDirection(row: number, col: number): string | null {
        return laneDirectionMap.value.get(`${row},${col}`) || null
    }

    function getLaneArrow(direction: string | null): string {
        const arrows: Record<string, string> = {
        W: '←',
        E: '→',
        N: '↑',
        S: '↓'
        }
        return direction ? arrows[direction] || '' : ''
    }

    function getLaneArrowClass(direction: string | null): string {
        const classes: Record<string, string> = {
        W: 'text-yellow-400/50',
        E: 'text-green-400/50',
        N: 'text-green-400/50',
        S: 'text-red-400/50'
        }
        return direction ? classes[direction] || '' : ''
    }

    // Check which directions are used in this config
    const activeDirections = computed(() => {
        const active: string[] = []
        for (const [dir, cells] of Object.entries(laneConfig.value)) {
        if (cells.length > 0) active.push(dir)
        }
        return active
    })

    return {
        getLaneDirection,
        getLaneArrow,
        getLaneArrowClass,
        activeDirections
    }
}
