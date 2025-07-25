import { FileWithPath } from '@mantine/dropzone'

export interface GetPlaylistsRequest {
  currentPage?: number
  pageSize?: number
  orderBy?: string[]
  searchBy?: string[]
}

export interface GetPlaylistRequest {
  id: string
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

// songs

export interface GetPlaylistSongsRequest {
  id: string
  currentPage?: number
  pageSize?: number
  orderBy?: string[]
  searchBy?: string[]
}

export interface AddSongsToPlaylistRequest {
  id: string
  songIds: string[]
  forceAdd?: boolean
}

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
