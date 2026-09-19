import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useRef,
  useState
} from 'react'
import OutlineView from '../types/enums/OutlineView.ts'
import LocalStorageKeys from '../../../../../../../types/enums/keys/LocalStorageKeys.ts'
import { useSongContext } from '../../../context/SongContext.tsx'
import useLocalStorage from '../../../../../../../hooks/useLocalStorage.ts'

interface SongOutlineContextValue {
  view: OutlineView
  showDetails: boolean
  setView: (view: OutlineView) => void
  setShowDetails: (showDetails: boolean) => void
}

const SongOutlineContext = createContext<SongOutlineContextValue | undefined>(undefined)

interface SongOutlineProviderProps {
  children: ReactNode
  initialView?: OutlineView
  initialShowDetails?: boolean
}

export function SongOutlineProvider({
  children,
  initialView = OutlineView.Sections,
  initialShowDetails = false
}: SongOutlineProviderProps) {
  const { songId } = useSongContext()
  const previousSongId = useRef(songId)
  const [view, setView] = useLocalStorage<OutlineView>({
    key: LocalStorageKeys.SongOutlineView,
    defaultValue: initialView,
    deserialize: (item) => {
      const parsed = JSON.parse(item)
      return parsed === OutlineView.Parts || parsed === OutlineView.Sections
        ? parsed
        : OutlineView.Sections
    }
  })
  const [showDetails, setShowDetails] = useState(initialShowDetails)

  useEffect(() => {
    if (previousSongId.current !== songId) {
      setShowDetails(false)
      previousSongId.current = songId
    }
  }, [songId])

  const value: SongOutlineContextValue = {
    view,
    showDetails,
    setView,
    setShowDetails
  }

  return <SongOutlineContext.Provider value={value}>{children}</SongOutlineContext.Provider>
}

export function useSongOutlineContext() {
  const context = useContext(SongOutlineContext)
  if (!context) throw new Error('useSongOutlineContext must be used within SongOutlineProvider')
  return context
}
