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

interface DrawersState {
  artist: ArtistDrawer
  album: AlbumDrawer
  song: SongDrawer
  playlist: PlaylistDrawer
}

const initialState: DrawersState = {
  song: {
    open: false
  },
  album: {
    open: false
  },
  artist: {
    open: false
  },
  playlist: {
    open: false
  }
}

export const drawersSlice = createSlice({
  name: 'drawers',
  initialState,
  reducers: {
    openArtistDrawer: (state, action: PayloadAction<string>) => {
      state.artist.artistId = action.payload
      state.artist.open = true
    },
    closeArtistDrawer: (state) => {
      state.artist.open = false
    },
    deleteArtistDrawer: (state) => {
      state.artist.open = false
      state.artist.artistId = undefined
    },

    openAlbumDrawer: (state, action: PayloadAction<string>) => {
      state.album.albumId = action.payload
      state.album.open = true
    },
    closeAlbumDrawer: (state) => {
      state.album.open = false
    },
    deleteAlbumDrawer: (state) => {
      state.album.open = false
      state.album.albumId = undefined
    },

    openSongDrawer: (state, action: PayloadAction<string>) => {
      state.song.songId = action.payload
      state.song.open = true
    },
    closeSongDrawer: (state) => {
      state.song.open = false
    },
    deleteSongDrawer: (state) => {
      state.song.open = false
      state.song.songId = undefined
    },

    openPlaylistDrawer: (state, action: PayloadAction<string>) => {
      state.playlist.playlistId = action.payload
      state.playlist.open = true
    },
    closePlaylistDrawer: (state) => {
      state.playlist.open = false
    },
    deletePlaylistDrawer: (state) => {
      state.playlist.open = false
      state.playlist.playlistId = undefined
    }
  }
})

export const {
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
} = drawersSlice.actions

export default drawersSlice.reducer
