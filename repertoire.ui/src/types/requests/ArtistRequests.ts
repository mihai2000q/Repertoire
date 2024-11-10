import { FileWithPath } from '@mantine/dropzone'

export interface GetArtistsRequest {
  currentPage?: number
  pageSize?: number
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
