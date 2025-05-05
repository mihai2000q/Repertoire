import { createContext, MutableRefObject, ReactNode, useContext, useRef } from 'react'

const ScrollRefContext = createContext(null)

export function ScrollRefProvider({ children }: { children: ReactNode }) {
  const scrollRef = useRef<HTMLDivElement>(null)

  return <ScrollRefContext.Provider value={scrollRef}>{children}</ScrollRefContext.Provider>
}

export const useScrollRef = () => useContext<MutableRefObject<HTMLDivElement>>(ScrollRefContext)
