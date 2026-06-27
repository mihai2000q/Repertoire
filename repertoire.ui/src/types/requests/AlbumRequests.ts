import { FileWithPath } from '@mantine/dropzone'

export interface GetAlbumsRequest {
  currentPage?: number
  pageSize?: number
  orderBy?: string[]
  searchBy?: string[]
}

export interface GetAlbumRequest {
  id: string
  songsOrderBy?: string[]
}

export interface CreateAlbumRequest {
  title: string
  releaseDate?: Date | string
  artistId?: string
  artistName?: string
}

export interface AddPerfectRehearsalsToAlbumsRequest {
  ids: string[]
}

export interface BulkDeleteAlbumsRequest {
  ids: string[]
  withSongs?: boolean
}

export interface SaveImageToAlbumRequest {
  image: FileWithPath
  id: string
}

export interface DeleteAlbumRequest {
  id: string
  withSongs?: boolean
}
