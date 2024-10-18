export interface GetSongsRequest {
  currentPage?: number
  pageSize?: number
}

export interface CreateSongRequest {
  title: string
  isRecorded?: boolean
}

export interface UpdateSongRequest {
  id: string
  title: string
  isRecorded?: boolean
}
