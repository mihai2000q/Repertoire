export interface ArtistFiltersMetadata {
  minBandMembersCount: number
  maxBandMembersCount: number

  minAlbumsCount: number
  maxAlbumsCount: number

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
