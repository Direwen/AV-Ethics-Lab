import { CELL_DEFINITIONS } from '~/constants/mapDefinitions'
import type { CellDefinition } from '~/types/simulation'

export function useCellDefinition() {
    const getCellDefinition = (code: number | string): CellDefinition => {
        const def = CELL_DEFINITIONS[String(code)]
        return def ?? { class: 'bg-red-500', name: '', isInteractive: false }
    }

    return { getCellDefinition }
}
