import { createContext, ReactNode, useContext, useEffect, useRef, useState } from 'react'
import DragSelect, { DSInputElement } from 'dragselect'
import { createStyles } from '@mantine/emotion'
import { alpha } from '@mantine/core'
import { useMain } from './MainContext.tsx'

const useStyles = createStyles((theme) => ({
  selector: {
    backgroundColor: `${alpha(theme.colors.gray[3], 0.33)} !important`,
    border: `2px solid ${alpha(theme.colors.gray[4], 0.8)} !important`,
    borderRadius: '12px'
  }
}))

interface DragSelectProviderProps {
  children: ReactNode
  settings?: ConstructorParameters<typeof DragSelect<DSInputElement>>[0]
  data?: unknown
}

interface DragSelectReturnType {
  dragSelect: DragSelect<DSInputElement> | undefined
  selectedIds: string[]
  clearSelection: () => void
}

const DragSelectContext = createContext<DragSelectReturnType>({
  dragSelect: undefined,
  selectedIds: [],
  clearSelection: () => undefined
})

export function DragSelectProvider({ children, data, settings = {} }: DragSelectProviderProps) {
  const [selectedIds, setSelectedIds] = useState<string[]>([])
  const dragSelectRef = useRef<DragSelect<DSInputElement>>(null)
  const { classes } = useStyles()
  const { ref: appRef } = useMain()

  // Helper: resolve area element from settings.area (supports HTMLElement, ref, or function)
  const resolveArea = (): HTMLElement | null => {
    const area = settings.area
    if (!area) return document.body
    if (area instanceof HTMLElement) return area
    if (area && 'current' in area && area.current instanceof HTMLElement) {
      return area.current
    }
    return document.body
  }

  // Initialize DragSelect instance
  useEffect(() => {
    const areaElement = resolveArea()
    if (!(areaElement instanceof HTMLElement)) return () => {}

    const fullSettings = {
      draggability: false,
      immediateDrag: false,
      keyboardDrag: false,
      multiSelectKeys: ['Control', 'Shift'],
      selectorClass: classes.selector,
      ...settings,
      area: areaElement
    }

    const ds = new DragSelect(fullSettings)
    dragSelectRef.current = ds

    return () => {
      ds.stop()
      dragSelectRef.current = undefined
    }
  }, [])

  // sync settings
  useEffect(() => {
    const ds = dragSelectRef.current
    if (!ds) return

    const { area: _area, ...updatableSettings } = settings
    ds.setSettings({
      draggability: false,
      immediateDrag: false,
      keyboardDrag: false,
      multiSelectKeys: ['Control', 'Shift'],
      selectorClass: classes.selector,
      ...updatableSettings
    })
    return
  }, [settings, classes.selector])

  // Selection changes
  useEffect(() => {
    const ds = dragSelectRef.current
    if (!ds) return () => {}

    const selectionChange = () => {
      const newIds = ds.getSelection().map((el) => el.id)
      setSelectedIds(newIds)
    }

    ds.subscribe('DS:start', selectionChange)
    ds.subscribe('DS:select', selectionChange)
    ds.subscribe('DS:unselect', selectionChange)

    return () => {
      ds.unsubscribe('DS:start', selectionChange)
      ds.unsubscribe('DS:select', selectionChange)
      ds.unsubscribe('DS:unselect', selectionChange)
    }
  }, [dragSelectRef.current])

  // Click outside detection
  useEffect(() => {
    const clickOutside = (event: MouseEvent) => {
      const areaElement = resolveArea()
      if (
        selectedIds.length > 0 &&
        appRef.current?.contains(event.target as Node) &&
        !areaElement.contains(event.target as Node)
      ) {
        dragSelectRef.current?.clearSelection()
      }
    }
    const appNode = appRef.current
    appNode?.addEventListener('click', clickOutside)
    return () => appNode?.removeEventListener('click', clickOutside)
  }, [appRef, selectedIds.length, settings.area])

  useEffect(() => handleClearSelection, [data])

  function handleClearSelection() {
    dragSelectRef.current?.clearSelection()
  }

  return (
    <DragSelectContext.Provider
      value={{
        dragSelect: dragSelectRef.current,
        selectedIds: selectedIds,
        clearSelection: handleClearSelection
      }}
    >
      {children}
    </DragSelectContext.Provider>
  )
}

export function useDragSelect() {
  return useContext(DragSelectContext)
}
