export interface AlbumFiltersMetadata {
  artistIds: string[]

  minReleaseDate?: string
  maxReleaseDate?: string

  minSongsCount: number
  maxSongsCount: number

  minRehearsals: number
  maxRehearsals: number

  minConfidence: number
  maxConfidence: number

  minProgress: number
  maxProgress: number

  minLastTimePlayed?: string
  maxLastTimePlayed?: string
}
