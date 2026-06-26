import { FileWithPath } from '@mantine/dropzone'
import Difficulty from '../enums/Difficulty.ts'

export interface GetSongsRequest {
  currentPage?: number
  pageSize?: number
  orderBy?: string[]
  searchBy?: string[]
  with?: string[]
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

export interface AddCustomSongRehearsalRequest {
  id: string
  arrangementId: string
}

export interface AddCustomSongRehearsalsRequest {
  requests: AddCustomSongRehearsalRequest[]
}

export interface AddPerfectSongRehearsalRequest {
  id: string
}

export interface AddPerfectSongRehearsalsRequest {
  ids: string[]
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
  albumId?: string
  artistId?: string
}

export interface UpdateSongSettingsRequest {
  settingsId: string
  defaultBandMemberId?: string
  defaultInstrumentId?: string
}

export interface BulkDeleteSongsRequest {
  ids: string[]
}

export interface SaveImageToSongRequest {
  image: FileWithPath
  id: string
}

export interface GetSongArrangementsRequest {
  songId: string
}
