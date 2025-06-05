export interface AddToPlaylistResponse {
  success: boolean
}

export interface AddArtistsToPlaylistResponse extends AddToPlaylistResponse {
  duplicateSongIds: string[]
  duplicateArtistIds: string[]
  addedSongIds: string[]
}

export interface AddAlbumsToPlaylistResponse extends AddToPlaylistResponse {
  duplicateSongIds: string[]
  duplicateAlbumIds: string[]
  addedSongIds: string[]
}

export interface AddSongsToPlaylistResponse extends AddToPlaylistResponse {
  duplicates: string[]
  added: string[]
}
