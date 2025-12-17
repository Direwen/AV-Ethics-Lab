import { CELL_DEFINITIONS } from '~/constants/mapDefinitions'
import type { CellDefinition } from '~/components/simulation/Cell.vue'

export function useCellDefinition() {
    const getCellDefinition = (code: number | string): CellDefinition => {
        const def = CELL_DEFINITIONS[String(code)]
        return def ? { class: def.class, name: def.name, allow: [] } : { class: 'bg-red-500', allow: [] }
    }

    return { getCellDefinition }
}
