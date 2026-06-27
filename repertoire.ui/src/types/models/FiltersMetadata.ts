import Difficulty from '../enums/Difficulty.ts'

export interface PlaylistFiltersMetadata {
  minSongsCount: number
  maxSongsCount: number
}

export interface SongFiltersMetadata {
  artistIds: string[]

  albumIds: string[]

  minReleaseDate?: string
  maxReleaseDate?: string

  minBpm?: number
  maxBpm?: number

  difficulties: Difficulty[]
  guitarTuningIds: string[]
  instrumentIds: string[]

  minSectionsCount: number
  maxSectionsCount: number

  minSolosCount: number
  maxSolosCount: number

  minRiffsCount: number
  maxRiffsCount: number

  minRehearsals: number
  maxRehearsals: number

  minConfidence: number
  maxConfidence: number

  minProgress: number
  maxProgress: number

  minLastTimePlayed?: string
  maxLastTimePlayed?: string
}
