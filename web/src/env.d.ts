/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module 'vuedraggable' {
  import { DefineComponent, VNode } from 'vue'

  interface DraggableProps {
    modelValue?: any[]
    itemKey?: string | ((item: any) => string | number)
    handle?: string
    animation?: number | string
    group?: string | { name: string; pull?: boolean | string; put?: boolean | string }
    disabled?: boolean
    ghostClass?: string
    dragClass?: string
    chosenClass?: string
    onStart?: (evt: any) => void
    onEnd?: (evt: any) => void
    onMove?: (evt: any) => boolean
    onChange?: (evt: any) => void
    'onUpdate:modelValue'?: (value: any[]) => void
  }

  const Draggable: DefineComponent<DraggableProps>
  export default Draggable
}
