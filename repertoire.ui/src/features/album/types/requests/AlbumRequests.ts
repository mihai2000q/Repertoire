export interface UpdateAlbumRequest {
  id: string
  title: string
  releaseDate?: Date | string
  artistId?: string
}

export interface AddSongsToAlbumRequest {
  id: string
  songIds: string[]
}

export interface MoveSongFromAlbumRequest {
  id: string
  songId: string
  overSongId: string
}

export interface RemoveSongsFromAlbumRequest {
  id: string
  songIds: string[]
}
