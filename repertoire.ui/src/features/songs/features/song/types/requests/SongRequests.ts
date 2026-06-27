import Difficulty from '../../../../../../types/enums/Difficulty.ts'

export interface UpdateSongRequest {
  id: string
  title: string
  description: string
  isRecorded?: boolean
  bpm?: number
  songsterrLink?: string
  youtubeLink?: string
  releaseDate?: Date | string
  difficulty?: Difficulty
  guitarTuningId?: string
  albumId?: string
  artistId?: string
}
