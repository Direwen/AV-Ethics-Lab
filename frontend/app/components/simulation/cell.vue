<template>
    <div
        ref="cellRef"
        class="relative w-4 h-4 sm:w-6 sm:h-6 md:w-8 md:h-8 lg:w-12 lg:h-12 flex items-center justify-center transition-all duration-30"
        :class="[
            definition?.class, 
            isInteractive ? 'hover:brightness-110 cursor-pointer' : '',
            showMenu ? 'ring-2 ring-cyan-400 z-50' : ''
        ]"
        :title="definition?.name" 
        role="gridcell"
        @click="handleClick"
        @mouseenter="isHovered = true"  
        @mouseleave="isHovered = false"
    >
        <EntityRenderer 
            v-if="entities?.length" 
            :entities="entities" 
            :parent-hover="isHovered"
        />

        <div 
            v-if="showMenu"
            class="absolute top-full mt-2 left-1/2 -translate-x-1/2 
                bg-slate-900 border border-slate-700 rounded-lg shadow-xl 
                z-[100] min-w-[160px] flex flex-col overflow-hidden"
        >
            <div class="bg-slate-800 px-3 py-1.5 text-[10px] uppercase font-bold text-slate-400 border-b border-slate-700">
                Select Entity
            </div>
            
            <button
                v-for="entity in entities"
                :key="entity.id"
                @click.stop="handleEntitySelect(entity)"
                class="flex items-center gap-3 px-3 py-2 hover:bg-cyan-500/20 hover:text-cyan-400 transition-colors text-left w-full group border-b border-slate-800 last:border-0"
            >
                <span class="text-xl group-hover:scale-110 transition-transform">{{ entity.emoji }}</span>
                <div class="flex flex-col leading-none">
                    <span class="text-xs font-medium text-slate-200">
                        {{ entity.metadata.name }}
                    </span>
                    <span v-if="entity.metadata.is_occluded" class="text-[9px] text-red-400 mt-0.5">
                        (Hidden)
                    </span>
                </div>
            </button>
        </div>

        <div 
            v-if="showMenu" 
            class="fixed inset-0 z-[90] cursor-default" 
            @click.stop="showMenu = false"
        ></div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import EntityRenderer from '~/components/simulation/EntityRenderer.vue'
import type { Entity } from '~/types/simulation'

export interface CellDefinition {
    name?: string
    class: string
    allow?: string[]
    risk_factor?: string
}

const props = defineProps<{
    cellCode: number | string
    definition: CellDefinition
    entities?: Entity[]
}>()

const emit = defineEmits(['cell-click', 'entity-select'])
const isHovered = ref(false)
const showMenu = ref(false)

const isInteractive = computed(() => {
    return (props.definition?.allow?.length ?? 0) > 0
})

function handleClick() {
    if (!isInteractive.value) return
    
    // CASE 1: Empty Cell -> Normal Click
    if (!props.entities || props.entities.length === 0) {
        emit('cell-click', {
            code: props.cellCode,
            allowed: props.definition.allow,
            type: props.definition.name
        })
        return
    }

    // CASE 2: Single Entity -> Auto Select
    if (props.entities.length === 1) {
        emit('entity-select', props.entities[0])
        return
    }

    // CASE 3: Stack -> Open Menu
    showMenu.value = !showMenu.value
}

function handleEntitySelect(entity: Entity) {
    emit('entity-select', entity)
    showMenu.value = false
}
</script>