import { FileWithPath } from '@mantine/dropzone'

export interface GetPlaylistsRequest {
  currentPage?: number
  pageSize?: number
  orderBy?: string[]
  searchBy?: string[]
}

export interface GetPlaylistRequest {
  id: string
  songsOrderBy?: string[]
}

export interface CreatePlaylistRequest {
  title: string
  description: string
}

export interface UpdatePlaylistRequest {
  id: string
  title: string
  description: string
}

export interface SaveImageToPlaylistRequest {
  image: FileWithPath
  id: string
}

export interface AddArtistsToPlaylistRequest {
  id: string
  artistIds: string[]
  forceAdd?: boolean
}

export interface AddAlbumsToPlaylistRequest {
  id: string
  albumIds: string[]
  forceAdd?: boolean
}

export interface AddSongsToPlaylistRequest {
  id: string
  songIds: string[]
  forceAdd?: boolean
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
