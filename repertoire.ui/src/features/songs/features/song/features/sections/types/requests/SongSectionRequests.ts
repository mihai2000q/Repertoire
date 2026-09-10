export interface CreateSongSectionRequest {
  songId: string
  typeId: string
  name: string
  partIds: string[]
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
}
