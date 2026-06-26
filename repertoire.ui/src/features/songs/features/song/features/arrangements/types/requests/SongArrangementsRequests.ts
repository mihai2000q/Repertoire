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
