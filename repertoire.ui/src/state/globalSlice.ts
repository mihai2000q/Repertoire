import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface SongDrawer {
  songId?: string
  open: boolean
}

interface AlbumDrawer {
  albumId?: string
  open: boolean
}

export interface GlobalState {
  userId?: string | undefined
  errorPath?: string | undefined
  songDrawer: SongDrawer
  albumDrawer: AlbumDrawer
}

const initialState: GlobalState = {
  songDrawer: {
    open: false
  },
  albumDrawer: {
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
    openSongDrawer: (state, action: PayloadAction<string | undefined>) => {
      state.songDrawer.songId = action.payload
      state.songDrawer.open = true
    },
    closeSongDrawer: (state) => {
      state.songDrawer.open = false
    },
    openAlbumDrawer: (state, action: PayloadAction<string | undefined>) => {
      state.albumDrawer.albumId = action.payload
      state.albumDrawer.open = true
    },
    closeAlbumDrawer: (state) => {
      state.albumDrawer.open = false
    }
  }
})

export const {
  setUserId,
  setErrorPath,
  openSongDrawer,
  closeSongDrawer,
  openAlbumDrawer,
  closeAlbumDrawer
} = globalSlice.actions

export default globalSlice.reducer
