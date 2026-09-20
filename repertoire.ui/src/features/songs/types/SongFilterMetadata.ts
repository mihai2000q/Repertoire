import Difficulty from '../../../types/enums/Difficulty.ts'

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

  minPartsCount: number
  maxPartsCount: number

  minSectionsCount: number
  maxSectionsCount: number

  minSolosCount: number
  maxSolosCount: number

  minRehearsals: number
  maxRehearsals: number

  minConfidence: number
  maxConfidence: number

  minProgress: number
  maxProgress: number

  minLastTimePlayed?: string
  maxLastTimePlayed?: string
}
