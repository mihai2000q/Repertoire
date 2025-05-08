import { createContext, MutableRefObject, ReactNode, useEffect, useRef, useState } from 'react'

export interface MainScrollContextReturnType {
  ref: MutableRefObject<HTMLDivElement>
  isTopScrollPositionOver0: boolean
}

export const MainScrollContext = createContext(null)

export function MainScrollProvider({ children }: { children: ReactNode }) {
  const ref = useRef<HTMLDivElement>(null)
  const [isTopScrollPositionOver0, setIsTopScrollPositionOver0] = useState(false)

  useEffect(() => {
    const handleScroll = () => {
      if (ref.current) {
        setIsTopScrollPositionOver0(ref.current.scrollTop > 0)
      }
    }

    const currentRef = ref.current
    currentRef?.addEventListener('scroll', handleScroll)

    return () => {
      currentRef?.removeEventListener('scroll', handleScroll)
    }
  }, [])

  return (
    <MainScrollContext.Provider value={{ ref, isTopScrollPositionOver0 }}>
      {children}
    </MainScrollContext.Provider>
  )
}
