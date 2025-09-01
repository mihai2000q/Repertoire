import { useEffect, useRef } from 'react'
import { useDragSelect } from '../context/DragSelectContext.tsx'

export default function useDragSelectSelectableRef<T extends HTMLElement>() {
  const ref = useRef<T>()
  const { dragSelect } = useDragSelect()
  useEffect(() => {
    if (!ref.current || !dragSelect) return
    dragSelect.addSelectables(ref.current)
    return () => {
      if (dragSelect && ref.current) dragSelect.removeSelectables(ref.current)
    }
  }, [ref, dragSelect])
  return ref
}
