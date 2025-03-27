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
  userID?: string | undefined
  documentTitle?: string | undefined
  errorPath?: string | undefined
  artistDrawer: ArtistDrawer
  albumDrawer: AlbumDrawer
  songDrawer: SongDrawer
}

const initialState: GlobalState = {
  userID: undefined,
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
    setUserID: (state, action: PayloadAction<string>) => {
      state.userID = action.payload
    },
    setDocumentTitle: (state, action: PayloadAction<string | undefined>) => {
      state.documentTitle = action.payload
    },
    setErrorPath: (state, action: PayloadAction<string | undefined>) => {
      state.errorPath = action.payload
    },
    openArtistDrawer: (state, action: PayloadAction<string>) => {
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

    openAlbumDrawer: (state, action: PayloadAction<string>) => {
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

    openSongDrawer: (state, action: PayloadAction<string>) => {
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
  setUserID,
  setDocumentTitle,
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
