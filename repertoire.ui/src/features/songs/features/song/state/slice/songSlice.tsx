import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import Song, { SongSettings } from '../../../../../../types/models/Song.ts'
import { BandMember } from '../../../../../../types/models/Artist.ts'

interface SongState {
  songId: string
  isArtistBand: boolean
  artistBandMembers?: BandMember[] | undefined
  settings?: SongSettings
  defaultArrangementId?: string | null
}

const initialState: SongState = {
  songId: "",
  isArtistBand: false,
  artistBandMembers: undefined,
}

export const songSlice = createSlice({
  name: 'song',
  initialState,
  reducers: {
    setSong: (state, action: PayloadAction<Song>) => {
      const song = action.payload
      const artist = song.artist

      state.songId = song.id
      state.settings = song.settings
      state.defaultArrangementId = song.defaultArrangementId
      state.isArtistBand = artist?.isBand
      state.artistBandMembers = artist?.isBand === false ? undefined : artist?.bandMembers
    }
  }
})

export const { setSong } = songSlice.actions

export default songSlice.reducer
