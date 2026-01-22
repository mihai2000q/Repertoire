import {
  createContext,
  ReactNode,
  RefObject,
  useContext,
  useEffect,
  useState
} from 'react'

interface MainContextReturnType {
  ref: RefObject<HTMLDivElement>
  mainScroll: {
    ref: RefObject<HTMLDivElement>
    isPositionOver0: boolean
  }
}

export const MainContext = createContext<MainContextReturnType>({
  ref: null,
  mainScroll: {
    ref: null,
    isPositionOver0: false
  }
})

interface MainContextProps {
  children: ReactNode
  appRef: RefObject<HTMLDivElement>
  scrollRef: RefObject<HTMLDivElement>
}

export function MainProvider({ children, appRef, scrollRef }: MainContextProps) {
  const [isTopScrollPositionOver0, setIsTopScrollPositionOver0] = useState(false)

  useEffect(() => {
    const handleScroll = () => {
      setIsTopScrollPositionOver0(scrollRef.current?.scrollTop > 0)
    }

    scrollRef.current?.addEventListener('scroll', handleScroll)
    return () => {
      scrollRef.current?.removeEventListener('scroll', handleScroll)
    }
  }, [])

  return (
    <MainContext.Provider
      value={{
        ref: appRef,
        mainScroll: {
          ref: scrollRef,
          isPositionOver0: isTopScrollPositionOver0
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
