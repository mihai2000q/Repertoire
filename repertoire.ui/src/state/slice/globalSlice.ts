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

interface PlaylistDrawer {
  playlistId?: string
  open: boolean
}

interface GlobalState {
  userId?: string | undefined
  documentTitle?: string | undefined
  errorPath?: string | undefined
  artistDrawer: ArtistDrawer
  albumDrawer: AlbumDrawer
  songDrawer: SongDrawer
  playlistDrawer: PlaylistDrawer
}

const initialState: GlobalState = {
  userId: undefined,
  songDrawer: {
    open: false
  },
  albumDrawer: {
    open: false
  },
  artistDrawer: {
    open: false
  },
  playlistDrawer: {
    open: false
  }
}

export const globalSlice = createSlice({
  name: 'global',
  initialState,
  reducers: {
    setUserId: (state, action: PayloadAction<string>) => {
      state.userId = action.payload
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
    },

    openPlaylistDrawer: (state, action: PayloadAction<string>) => {
      state.playlistDrawer.playlistId = action.payload
      state.playlistDrawer.open = true
    },
    closePlaylistDrawer: (state) => {
      state.playlistDrawer.open = false
    },
    deletePlaylistDrawer: (state) => {
      state.playlistDrawer.open = false
      state.playlistDrawer.playlistId = undefined
    }
  }
})

export const {
  setUserId,
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
  deleteSongDrawer,
  openPlaylistDrawer,
  closePlaylistDrawer,
  deletePlaylistDrawer
} = globalSlice.actions

export default globalSlice.reducer
