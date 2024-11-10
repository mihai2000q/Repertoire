import {FileWithPath} from "@mantine/dropzone";

export interface GetAlbumsRequest {
  currentPage?: number
  pageSize?: number
}

export interface CreateAlbumRequest {
  title: string
  releaseDate: string
  artistId?: string
  artistName?: string
}

export interface UpdateAlbumRequest {
  id: string
  title: string
  releaseDate: string
}

export interface SaveImageToAlbumRequest {
  image: FileWithPath
  id: string
}
