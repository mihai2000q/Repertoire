import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface ArtistDrawer {
  artistId?: string
  open: boolean
}

interface AlbumDrawer {
  albumId?: string
  open: boolean
}

interface SongDrawer {
  songId?: string
  open: boolean
}

export interface GlobalState {
  userId?: string | undefined
  errorPath?: string | undefined
  artistDrawer: ArtistDrawer
  albumDrawer: AlbumDrawer
  songDrawer: SongDrawer
}

const initialState: GlobalState = {
  songDrawer: {
    open: false
  },
  albumDrawer: {
    open: false
  },
  artistDrawer: {
    open: false
  }
}

export const globalSlice = createSlice({
  name: 'global',
  initialState,
  reducers: {
    setUserId: (state, action: PayloadAction<string | undefined>) => {
      state.userId = action.payload
    },
    setErrorPath: (state, action: PayloadAction<string | undefined>) => {
      state.errorPath = action.payload
    },
    openArtistDrawer: (state, action: PayloadAction<string | undefined>) => {
      state.artistDrawer.artistId = action.payload
      state.artistDrawer.open = true
    },
    closeArtistDrawer: (state) => {
      state.artistDrawer.open = false
    },
    openAlbumDrawer: (state, action: PayloadAction<string | undefined>) => {
      state.albumDrawer.albumId = action.payload
      state.albumDrawer.open = true
    },
    closeAlbumDrawer: (state) => {
      state.albumDrawer.open = false
    },
    openSongDrawer: (state, action: PayloadAction<string | undefined>) => {
      state.songDrawer.songId = action.payload
      state.songDrawer.open = true
    },
    closeSongDrawer: (state) => {
      state.songDrawer.open = false
    }
  }
})

export const {
  setUserId,
  setErrorPath,
  openArtistDrawer,
  closeArtistDrawer,
  openAlbumDrawer,
  closeAlbumDrawer,
  openSongDrawer,
  closeSongDrawer
} = globalSlice.actions

export default globalSlice.reducer
