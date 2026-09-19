export interface GetSongPartsRequest {
  songId: string
}

export interface CreateSongPartRequest {
  songId: string
  name: string
  instrumentId?: string
  sectionId?: string
  bandMemberId?: string
}

export interface BulkUpdateSongPartsRequest {
  requests: { id: string; rehearsals: number, confidence: number }[]
  songId: string
}

export interface UpdateSongPartRequest {
  id: string
  name: string
  rehearsals: number
  confidence: number
  instrumentId?: string
  sectionId?: string
  bandMemberId?: string
}

export interface UpdateAllSongPartsRequest {
  songId: string
  bandMemberId?: string
  instrumentId?: string
}

export interface MoveSongPartInSongRequest {
  id: string
  overId: string
  songId: string
}

export interface BulkDeleteSongPartsRequest {
  ids: string[]
  songId: string
}

export interface DeleteSongPartRequest {
  id: string
  songId: string
}
