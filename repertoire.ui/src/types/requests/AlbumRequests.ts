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

export interface UpdateAlbumRequest {
  id: string
  title: string
  releaseDate?: Date | string
  artistId?: string
}

export interface SaveImageToAlbumRequest {
  image: FileWithPath
  id: string
}

export interface DeleteAlbumRequest {
  id: string
  withSongs?: boolean
}

export interface AddSongsToAlbumRequest {
  id: string
  songIds: string[]
}

export interface MoveSongFromAlbumRequest {
  id: string
  songId: string
  overSongId: string
}

export interface RemoveSongsFromAlbumRequest {
  id: string
  songIds: string[]
}
