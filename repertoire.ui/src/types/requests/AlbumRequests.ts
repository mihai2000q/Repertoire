import { FileWithPath } from '@mantine/dropzone'

export interface GetAlbumsRequest {
  currentPage?: number
  pageSize?: number
  orderBy?: string[]
}

export interface CreateAlbumRequest {
  title: string
  releaseDate?: Date
  artistId?: string
  artistName?: string
}

export interface UpdateAlbumRequest {
  id: string
  title: string
  releaseDate?: Date
}

export interface SaveImageToAlbumRequest {
  image: FileWithPath
  id: string
}

export interface AddSongsToAlbumRequest {
  id: string
  songIds: string[]
}

export interface RemoveSongsFromAlbumRequest {
  id: string
  songIds: string[]
}
