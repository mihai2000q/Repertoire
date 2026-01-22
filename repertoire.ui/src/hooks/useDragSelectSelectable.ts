import { RefObject, useEffect, useRef, useState } from 'react'
import { useDragSelect } from '../context/DragSelectContext.tsx'

interface UseDragSelectSelectableReturnType<T> {
  ref: RefObject<T>
  isDragSelected: boolean
  isDragSelecting: boolean
}

export default function useDragSelectSelectable<T extends HTMLElement>(
  id: string
): UseDragSelectSelectableReturnType<T> {
  const ref = useRef<T>(null)
  const { dragSelect, selectedIds } = useDragSelect()
  const [isDragSelected, setIsDragSelected] = useState(false)
  const [isDragging, setIsDragging] = useState(false)

  useEffect(() => {
    if (!ref.current || !dragSelect) return () => {}
    dragSelect.addSelectables(ref.current)
    return () => {
      if (dragSelect && ref.current) dragSelect.removeSelectables(ref.current)
    }
  }, [ref, dragSelect])

  useEffect(() => {
    setIsDragSelected(selectedIds.some((i) => i === id))
    setIsDragging(selectedIds.length > 0)
  }, [selectedIds])

  return {
    ref: ref,
    isDragSelected: isDragSelected,
    isDragSelecting: isDragging
  }
}
