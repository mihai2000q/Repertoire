import { createContext, ReactNode, useContext } from 'react'
import Song, { SongSettings } from '../../../../../types/models/Song.ts'
import { BandMember } from '../../../../../types/models/Artist.ts'

interface SongContextValue {
  songId: string
  isArtistBand: boolean
  artistBandMembers?: BandMember[]
  settings: SongSettings
  defaultArrangementId?: string
}

const defaultSongContext: SongContextValue = {
  songId: '',
  isArtistBand: false,
  artistBandMembers: undefined,
  settings: undefined,
  defaultArrangementId: undefined
}

const SongContext = createContext<SongContextValue>(defaultSongContext)

interface SongProviderProps {
  children: ReactNode
  song: Song
}

export function SongProvider({ children, song }: SongProviderProps) {
  const artist = song.artist

  return (
    <SongContext.Provider
      value={{
        songId: song.id,
        settings: song.settings,
        isArtistBand: artist?.isBand ?? false,
        artistBandMembers: artist?.isBand === false ? undefined : artist?.bandMembers,
        defaultArrangementId: song.defaultArrangementId
      }}
    >
      {children}
    </SongContext.Provider>
  )
}

export function useSongContext() {
  return useContext(SongContext)
}
