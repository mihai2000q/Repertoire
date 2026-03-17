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

// Sections

export interface CreateSongSectionRequest {
  songId: string
  typeId: string
  name: string
  instrumentId?: string
  bandMemberId?: string
}

export interface BulkRehearsalsSongSectionsRequest {
  sections: { id: string; rehearsals: number }[]
  songId: string
}

export interface UpdateSongSectionRequest {
  id: string
  typeId: string
  name: string
  rehearsals: number
  confidence: number
  bandMemberId?: string
  instrumentId?: string
}

export interface UpdateAllSongSectionsRequest {
  songId: string
  bandMemberId?: string
  instrumentId?: string
}

export interface MoveSongSectionRequest {
  id: string
  overId: string
  songId: string
}

export interface BulkDeleteSongSectionsRequest {
  ids: string[]
  songId: string
}

export interface DeleteSongSectionRequest {
  id: string
  songId: string
}

// Arrangements

export interface GetSongArrangementsRequest {
  songId: string
}

export interface CreateSongArrangementRequest {
  songId: string
  name: string
}

export interface BulkUpdateSongArrangementsRequest {
  songId: string
  requests: UpdateSongArrangementRequest[]
}

export interface UpdateSongArrangementRequest {
  id: string
  name: string
  occurrences: UpdateSongSectionOccurrencesRequest[]
}

export interface UpdateSongSectionOccurrencesRequest {
  sectionId: string
  occurrences: number
}

export interface UpdateDefaultSongArrangementRequest {
  id: string | null
  songId: string
}

export interface MoveSongArrangementRequest {
  id: string
  overId: string
  songId: string
}

export interface DeleteSongArrangementRequest {
  id: string
  songId: string
}
