import { RefObject, useEffect, useRef, useState } from 'react'
import { useClickSelect } from '../context/ClickSelectContext.tsx'

interface UseClickSelectSelectableReturnType<T> {
  ref: RefObject<T>
  isClickSelected: boolean
  isClickSelectionActive: boolean
  isLastInSelection: boolean
}

export default function useClickSelectSelectable<T extends HTMLElement>(
  id: string
): UseClickSelectSelectableReturnType<T> {
  const ref = useRef<T>()
  const { addSelectable, removeSelectable, selectedIds, selectables } = useClickSelect()
  const [isSelected, setIsSelected] = useState(false)
  const [isActive, setIsActive] = useState(false)
  const [isLastInSelection, setIsLastInSelection] = useState(false)

  useEffect(() => {
    if (!ref.current) return () => {}
    addSelectable(id, ref.current)
    return () => removeSelectable(id)
  }, [ref])

  useEffect(() => {
    const isSelected = selectedIds.some((i) => i === id)
    setIsSelected(isSelected)
    setIsActive(selectedIds.length > 0)

    if (!isSelected || selectedIds.length === 0) {
      setIsLastInSelection(false)
      return
    }

    const index = selectables.findIndex((s) => s.id === id)
    setIsLastInSelection(index === selectables.length - 1 || !selectables[index + 1].selected)
  }, [selectedIds])

  return {
    ref: ref,
    isClickSelected: isSelected,
    isClickSelectionActive: isActive,
    isLastInSelection: isLastInSelection
  }
}
