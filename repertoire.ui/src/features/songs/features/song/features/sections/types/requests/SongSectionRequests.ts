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
