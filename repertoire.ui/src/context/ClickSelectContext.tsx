import React, { createContext, useContext, useEffect, useRef, useState } from 'react'
import { useMain } from './MainContext.tsx'

interface ClickSelectProviderProps {
  children: React.ReactNode
  data?: unknown
}

interface ClickSelectReturnType {
  selectables: ClickSelectable[]
  selectedIds: string[]
  isClickSelectionActive: boolean
  addSelectable: (id: string, el: HTMLElement) => void
  removeSelectable: (id: string) => void
  clearSelection: () => void
}

const ClickSelectContext = createContext<ClickSelectReturnType>({
  selectables: [],
  selectedIds: [],
  isClickSelectionActive: false,
  addSelectable: () => undefined,
  removeSelectable: () => undefined,
  clearSelection: () => undefined
})

interface ClickSelectable {
  id: string
  selected: boolean
}

export function ClickSelectProvider({ children, data }: ClickSelectProviderProps) {
  const selectables = useRef<ClickSelectable[]>([])
  const lastSelectedId = useRef('')
  const [selectedIds, setSelectedIds] = useState<string[]>([])
  const [isSelectionActive, setIsSelectionActive] = useState(false)
  const areaRef = useRef<HTMLSpanElement>(null)
  const { ref: appRef } = useMain()

  useEffect(() => handleClearSelection(), [data])

  useEffect(() => setIsSelectionActive(selectedIds.length > 0), [selectedIds])

  // Event delegation for clicks on selectable elements
  useEffect(() => {
    const container = areaRef.current
    if (!container) return

    const handleClick = (ev: MouseEvent) => {
      // Find the closest element with data-selectable-id
      let target = ev.target as HTMLElement | null
      while (target && target !== container) {
        const id = target.getAttribute('data-selectable-id')
        if (id) {
          if (ev.ctrlKey) ctrlClick(id)
          if (ev.shiftKey) shiftClick(id)
          return
        }
        target = target.parentElement
      }
    }

    container.addEventListener('click', handleClick)
    return () => container.removeEventListener('click', handleClick)
  }, [])

  // Click outside detection
  useEffect(() => {
    const clickOutside = (event: MouseEvent) => {
      if (
        isSelectionActive &&
        !areaRef.current?.contains(event.target as Node) &&
        appRef.current?.contains(event.target as Node)
      )
        handleClearSelection()
    }

    const appNode = appRef.current
    appNode?.addEventListener('click', clickOutside)
    return () => appNode?.removeEventListener('click', clickOutside)
  }, [isSelectionActive, appRef])

  function setNewIds(id: string) {
    const newIds = selectables.current.filter((s) => s.selected).map((s) => s.id)
    setSelectedIds(newIds)
    lastSelectedId.current = id
    if (newIds.length === 0) resetLastSelectedId()
  }
  function resetLastSelectedId() {
    lastSelectedId.current = selectables.current.length === 0 ? '' : selectables.current[0].id
  }

  function ctrlClick(id: string) {
    selectables.current = selectables.current.map((s) =>
      s.id !== id ? s : { ...s, selected: !s.selected }
    )
    setNewIds(id)
  }

  function shiftClick(id: string) {
    const currIndex = selectables.current.findIndex((s) => s.id === id)
    const currState = selectables.current[currIndex].selected
    const lastIndex = selectables.current.findIndex((s) => s.id === lastSelectedId.current)

    const indexes = lastIndex < currIndex ? [lastIndex, currIndex] : [currIndex, lastIndex]
    for (let i = indexes[0]; i <= indexes[1]; i++) {
      selectables.current[i].selected = !currState
    }

    setNewIds(id)
  }

  function handleAddSelectable(id: string, el: HTMLElement) {
    // Avoid duplicates
    if (selectables.current.some((s) => s.id === id)) return
    selectables.current.push({ id, selected: false })
    resetLastSelectedId()
    // Mark the element with a data attribute for event delegation
    el.setAttribute('data-selectable-id', id)
  }

  function handleRemoveSelectable(id: string) {
    selectables.current = selectables.current.filter((s) => s.id !== id)
    setSelectedIds((prev) => prev.filter((i) => i !== id))
    resetLastSelectedId()
    // Optionally remove the data attribute from the element, but not strictly necessary
  }

  function handleClearSelection() {
    selectables.current = selectables.current.map((s) => ({ ...s, selected: false }))
    setSelectedIds([])
    resetLastSelectedId()
  }

  return (
    <ClickSelectContext.Provider
      value={{
        selectables: selectables.current,
        selectedIds,
        isClickSelectionActive: isSelectionActive,
        addSelectable: handleAddSelectable,
        removeSelectable: handleRemoveSelectable,
        clearSelection: handleClearSelection
      }}
    >
      <span ref={areaRef} style={{ display: 'contents' }}>
        {children}
      </span>
    </ClickSelectContext.Provider>
  )
}

export function useClickSelect() {
  return useContext(ClickSelectContext)
}
