import { FileWithPath } from '@mantine/dropzone'
import Difficulty from '../../utils/enums/Difficulty.ts'

export interface GetSongsRequest {
  currentPage?: number
  pageSize?: number
  orderBy?: string[]
  searchBy?: string[]
}

export interface CreateSongRequest {
  title: string
  description: string
  bpm?: number
  releaseDate?: Date | string
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
  description: string
  isRecorded?: boolean
  bpm?: number
  songsterrLink?: string
  youtubeLink?: string
  releaseDate?: Date | string
  difficulty?: Difficulty
  guitarTuningId?: string
}

export interface SaveImageToSongRequest {
  image: FileWithPath
  id: string
}

export interface CreateSongSectionRequest {
  songId: string
  typeId: string
  name: string
}

export interface UpdateSongSectionRequest {
  id: string
  typeId: string
  name: string
  rehearsals: number
}

export interface MoveSongSectionRequest {
  id: string
  overId: string
  songId: string
}

export interface DeleteSongSectionRequest {
  id: string
  songId: string
}
