import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState
} from 'react'
import OutlineView from '../types/enums/OutlineView.ts'
import LocalStorageKeys from '../../../../../types/enums/keys/LocalStorageKeys.ts'
import { useSongContext } from '../../../context/SongContext.tsx'
import useLocalStorage from '../../../../../hooks/useLocalStorage.ts'

interface SongOutlineContextValue {
  view: OutlineView
  setView: (view: OutlineView) => void
  isDetailsShowing: (id: string) => boolean
  toggleDetails: (id: string) => void
  toggleAllDetails: (ids: string[]) => void
  areAllDetailsShowing: (ids: string[]) => boolean
}

const SongOutlineContext = createContext<SongOutlineContextValue | undefined>(undefined)

interface SongOutlineProviderProps {
  children: ReactNode
  initialView?: OutlineView
  initialDetailsSectionIds?: Set<string>
  initialDetailsPartIds?: Set<string>
}

export function SongOutlineProvider({
  children,
  initialView = OutlineView.Sections,
  initialDetailsSectionIds = new Set(),
  initialDetailsPartIds = new Set()
}: SongOutlineProviderProps) {
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

  const [openSectionIds, setOpenSectionIds] = useState(initialDetailsSectionIds)
  const [openPartIds, setOpenPartIds] = useState(initialDetailsPartIds)

  const openIds = view === OutlineView.Sections ? openSectionIds : openPartIds
  const setOpenIds = view === OutlineView.Sections ? setOpenSectionIds : setOpenPartIds

  const isDetailsShowing = useCallback((id: string) => openIds.has(id), [openIds])
  const toggleDetails = useCallback(
    (id: string) =>
      setOpenIds((prev) => {
        const next = new Set(prev)
        if (!next.delete(id)) next.add(id)
        return next
      }),
    [setOpenIds]
  )

  const areAllDetailsShowing = useCallback(
    (ids: string[]) => ids.length > 0 && ids.every((id) => openIds.has(id)),
    [openIds]
  )
  const toggleAllDetails = useCallback(
    (ids: string[]) =>
      setOpenIds((prev) => {
        const show = !(ids.length > 0 && ids.every((id) => prev.has(id)))
        const next = new Set(prev)
        ids.forEach((id) => (show ? next.add(id) : next.delete(id)))
        return next
      }),
    [setOpenIds]
  )

  // reset on song change
  const { songId } = useSongContext()
  const previousSongId = useRef(songId)
  useEffect(() => {
    if (previousSongId.current !== songId) {
      setOpenSectionIds(new Set())
      setOpenPartIds(new Set())
      previousSongId.current = songId
    }
  }, [songId])

  const value: SongOutlineContextValue = {
    view,
    setView,
    isDetailsShowing,
    toggleDetails,
    toggleAllDetails,
    areAllDetailsShowing
  }

  return <SongOutlineContext.Provider value={value}>{children}</SongOutlineContext.Provider>
}

export function useSongOutlineContext() {
  const context = useContext(SongOutlineContext)
  if (!context) throw new Error('useSongOutlineContext must be used within SongOutlineProvider')
  return context
}
