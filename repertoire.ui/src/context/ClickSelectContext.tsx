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
  const areaRef = useRef<HTMLSpanElement>(undefined)
  const { ref: appRef } = useMain()

  useEffect(() => handleClearSelection(), [data])

  useEffect(() => setIsSelectionActive(selectedIds.length > 0), [selectedIds])

  useEffect(() => {
    const clickOutside = (event: PointerEvent) => {
      if (
        isSelectionActive &&
        !areaRef.current?.contains(event.target as Node) &&
        appRef.current?.contains(event.target as Node)
      )
        handleClearSelection()
    }

    console.log(appRef)
    console.log(areaRef)
    appRef.current?.addEventListener('click', clickOutside)
    return () => appRef.current?.removeEventListener('click', clickOutside)
  }, [appRef, areaRef])

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

  // handlers
  function handleAddSelectable(id: string, el: HTMLElement) {
    selectables.current.push({ id: id, selected: false })
    resetLastSelectedId()

    el.onclick = (ev) => {
      if (ev.ctrlKey) ctrlClick(id)
      if (ev.shiftKey) shiftClick(id)
    }
  }

  function handleRemoveSelectable(id: string) {
    selectables.current = selectables.current.filter((s) => s.id !== id)
    setSelectedIds((prevState) => [...prevState.filter((id) => id !== id)])
    resetLastSelectedId()
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
        selectedIds: selectedIds,
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
