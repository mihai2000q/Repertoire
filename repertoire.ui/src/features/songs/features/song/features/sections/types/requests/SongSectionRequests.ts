export interface GetSongSectionsRequest {
  songId: string
}

export interface CreateSongSectionRequest {
  songId: string
  typeId: string
  name: string
  parts: CreateNewSongSectionPartRequest[]
}

interface CreateNewSongSectionPartRequest {
  bandMemberId: string
  partId: string
  newPart: { name: string; instrumentId: string }
}

export interface UpdateSongSectionRequest {
  id: string
  typeId: string
  name: string
  partIds: string[]
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
  withParts?: boolean
}
