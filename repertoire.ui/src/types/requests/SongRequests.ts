import { FileWithPath } from '@mantine/dropzone'

export interface GetSongsRequest {
  currentPage?: number
  pageSize?: number
}

export interface CreateSongRequest {
  title: string
  description: string
  bpm?: number
  releaseDate?: Date
  difficulty?: string
  songsterrLink?: string
  youtubeLink?: string

  sections?: CreateSectionRequest[]
  guitarTuningId?: string
  albumId?: string
  albumTitle?: string
  artistId?: string
  artistName?: string
}

export interface CreateSectionRequest {
  name: string
  typeId: string
}

export interface UpdateSongRequest {
  id: string
  title: string
  isRecorded?: boolean
}

export interface SaveImageToSongRequest {
  image: FileWithPath
  id: string
}
