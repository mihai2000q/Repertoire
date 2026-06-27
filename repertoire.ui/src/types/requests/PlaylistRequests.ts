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

export interface AddPerfectRehearsalsToPlaylistsRequest {
  ids: string[]
}

export interface BulkDeletePlaylistsRequest {
  ids: string[]
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

export interface AddPerfectPlaylistSongRehearsalsRequest {
  ids: string[]
}
