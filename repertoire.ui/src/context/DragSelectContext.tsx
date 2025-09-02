import React, { createContext, useContext, useEffect, useState } from 'react'
import DragSelect, { DSInputElement } from 'dragselect'
import { createStyles } from '@mantine/emotion'
import { alpha } from '@mantine/core'

const useStyles = createStyles((theme) => ({
  selector: {
    backgroundColor: `${alpha(theme.colors.gray[3], 0.33)} !important`,
    border: `2px solid ${alpha(theme.colors.gray[4], 0.8)} !important`,
    borderRadius: '12px'
  }
}))

interface DragSelectProviderProps {
  children: React.ReactNode
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
  clearSelection: () => undefined,
})

export function DragSelectProvider({ children, data, settings = {} }: DragSelectProviderProps) {
  const [selectedIds, setSelectedIds] = useState<string[]>([])
  const [dragSelect, setDragSelect] = useState<DragSelect<DSInputElement>>()
  const { classes } = useStyles()

  useEffect(() => {
    setDragSelect((prevState) => {
      if (prevState) return prevState
      return new DragSelect({})
    })

    return () => {
      if (dragSelect) {
        dragSelect.stop()
        setDragSelect(undefined)
      }
    }
  }, [dragSelect])

  useEffect(() => {
    const handleSelectionChange = () => {
      if (dragSelect?.getSelection().map((el) => el.id) !== selectedIds)
        setSelectedIds(dragSelect?.getSelection().map((el) => el.id))
    }
    dragSelect?.subscribe('DS:start', handleSelectionChange)
    dragSelect?.subscribe('DS:select', handleSelectionChange)
    dragSelect?.subscribe('DS:unselect', handleSelectionChange)
    return () => {
      dragSelect?.unsubscribe('DS:start', handleSelectionChange)
      dragSelect?.unsubscribe('DS:select', handleSelectionChange)
      dragSelect?.unsubscribe('DS:unselect', handleSelectionChange)
    }
  }, [dragSelect])

  useEffect(
    () =>
      dragSelect?.setSettings({
        draggability: false,
        immediateDrag: false,
        keyboardDrag: false,
        multiSelectKeys: ['Control'],
        selectorClass: classes.selector,
        ...settings
      }),
    [dragSelect, settings]
  )

  useEffect(() => handleClearSelection(), [data])

  function handleClearSelection() {
    dragSelect?.clearSelection()
  }

  return (
    <DragSelectContext.Provider
      value={{
        dragSelect: dragSelect,
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
