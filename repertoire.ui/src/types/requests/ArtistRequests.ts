import { FileWithPath } from '@mantine/dropzone'

export interface GetArtistsRequest {
  currentPage?: number
  pageSize?: number
  orderBy?: string[]
  searchBy?: string[]
}

export interface CreateArtistRequest {
  name: string
}

export interface UpdateArtistRequest {
  id: string
  name: string
}

export interface SaveImageToArtistRequest {
  image: FileWithPath
  id: string
}

export interface AddAlbumsToArtistRequest {
  id: string
  albumIds: string[]
}

export interface RemoveAlbumsFromArtistRequest {
  id: string
  albumIds: string[]
}

export interface AddSongsToArtistRequest {
  id: string
  songIds: string[]
}

export interface RemoveSongsFromArtistRequest {
  id: string
  songIds: string[]
}
