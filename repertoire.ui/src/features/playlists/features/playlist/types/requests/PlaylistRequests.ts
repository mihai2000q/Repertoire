export interface UpdatePlaylistRequest {
  id: string
  title: string
  description: string
}

// songs

export interface ShufflePlaylistSongsRequest {
  id: string
}

export interface MoveSongFromPlaylistRequest {
  id: string
  playlistSongId: string
  overPlaylistSongId: string
}

export interface RemoveSongsFromPlaylistRequest {
  id: string
  playlistSongIds: string[]
}
