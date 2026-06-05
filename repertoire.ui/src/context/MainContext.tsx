import { createContext, ReactNode, RefObject, useContext, useEffect, useRef } from 'react'

interface MainContextReturnType {
  ref: RefObject<HTMLDivElement>
  mainScroll: {
    ref: RefObject<HTMLDivElement>
    topbarRef: RefObject<HTMLDivElement>
  }
}

export const MainContext = createContext<MainContextReturnType>({
  ref: null,
  mainScroll: {
    ref: null,
    topbarRef: null
  }
})

interface MainContextProps {
  children: ReactNode
  appRef: RefObject<HTMLDivElement>
  scrollRef: RefObject<HTMLDivElement>
}

export function MainProvider({ children, appRef, scrollRef }: MainContextProps) {
  const topbarRef = useRef(null)

  useEffect(() => {
    const scrollEl = scrollRef.current
    const topbarEl = topbarRef.current
    if (!scrollEl || !topbarEl) return () => {}

    const handleScroll = () => {
      const hasScrolled = scrollEl.scrollTop > 0
      if (hasScrolled) topbarEl.classList.add('scrolled')
      else topbarEl.classList.remove('scrolled')
    }

    scrollEl.addEventListener('scroll', handleScroll)
    return () => scrollEl.removeEventListener('scroll', handleScroll)
  }, [])

  return (
    <MainContext.Provider
      value={{
        ref: appRef,
        mainScroll: {
          ref: scrollRef,
          topbarRef: topbarRef
        }
      }}
    >
      {children}
    </MainContext.Provider>
  )
}

export function useMain() {
  return useContext(MainContext)
}
