import { FileWithPath } from '@mantine/dropzone'

export interface GetArtistsRequest {
  currentPage?: number
  pageSize?: number
  orderBy?: string[]
  searchBy?: string[]
}

export interface AddPerfectRehearsalsToArtistsRequest {
  ids: string[]
}

export interface BulkDeleteArtistsRequest {
  ids: string[]
  withAlbums?: boolean
  withSongs?: boolean
}

export interface SaveImageToArtistRequest {
  image: FileWithPath
  id: string
}

export interface DeleteArtistRequest {
  id: string
  withAlbums?: boolean
  withSongs?: boolean
}
