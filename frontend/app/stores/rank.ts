import { defineStore } from 'pinia'
import type { Entity } from '~/types/simulation'

export const useRankStore = defineStore('rank', () => {
    const selectedEntity = ref<Entity | null>(null)

    function select(entity: Entity) {
        selectedEntity.value = entity
    }

    function clearSelection() {
        selectedEntity.value = null
    }

    function isSelected(entityId: string): boolean {
        return selectedEntity.value?.id === entityId
    }

    return {
        // State
        selectedEntity,

        // Actions
        select,
        clearSelection,
        isSelected
    }
})
