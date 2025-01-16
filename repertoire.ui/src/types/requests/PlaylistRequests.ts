import { FileWithPath } from '@mantine/dropzone'

export interface GetPlaylistsRequest {
  currentPage?: number
  pageSize?: number
  orderBy?: string[]
  searchBy?: string[]
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

export interface AddSongsToPlaylistRequest {
  id: string
  songIds: string[]
}

export interface MoveSongFromPlaylistRequest {
  id: string
  songId: string
}

export interface RemoveSongsFromPlaylistRequest {
  id: string
  songIds: string[]
}
