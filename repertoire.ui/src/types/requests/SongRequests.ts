import { FileWithPath } from '@mantine/dropzone'
import Difficulty from '../enums/Difficulty.ts'

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

export interface AddPerfectSongRehearsalRequest {
  id: string
}

export interface AddPerfectSongRehearsalsRequest {
  ids: string[]
}

export interface AddPartialSongRehearsalRequest {
  id: string
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

export interface UpdateSongSectionRequest {
  id: string
  typeId: string
  name: string
  rehearsals: number
  confidence: number
  bandMemberId?: string
  instrumentId?: string
}

export interface UpdateSongSectionsOccurrencesRequest {
  songId: string
  sections: { id: string, occurrences: number }[]
}

export interface UpdateSongSectionsPartialOccurrencesRequest {
  songId: string
  sections: { id: string, partialOccurrences: number }[]
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

export interface DeleteSongSectionRequest {
  id: string
  songId: string
}
