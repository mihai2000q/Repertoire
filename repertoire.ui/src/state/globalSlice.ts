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
    deleteArtistDrawer: (state) => {
      state.artistDrawer.open = false
      state.artistDrawer.artistId = undefined
    },

    openAlbumDrawer: (state, action: PayloadAction<string | undefined>) => {
      state.albumDrawer.albumId = action.payload
      state.albumDrawer.open = true
    },
    closeAlbumDrawer: (state) => {
      state.albumDrawer.open = false
    },
    deleteAlbumDrawer: (state) => {
      state.albumDrawer.open = false
      state.albumDrawer.albumId = undefined
    },

    openSongDrawer: (state, action: PayloadAction<string | undefined>) => {
      state.songDrawer.songId = action.payload
      state.songDrawer.open = true
    },
    closeSongDrawer: (state) => {
      state.songDrawer.open = false
    },
    deleteSongDrawer: (state) => {
      state.songDrawer.open = false
      state.songDrawer.songId = undefined
    }
  }
})

export const {
  setErrorPath,
  openArtistDrawer,
  closeArtistDrawer,
  deleteArtistDrawer,
  openAlbumDrawer,
  closeAlbumDrawer,
  deleteAlbumDrawer,
  openSongDrawer,
  closeSongDrawer,
  deleteSongDrawer
} = globalSlice.actions

export default globalSlice.reducer
