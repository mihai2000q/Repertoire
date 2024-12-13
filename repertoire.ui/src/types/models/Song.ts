import Difficulty from '../../utils/enums/Difficulty.ts'
import Album from './Album.ts'
import Artist from './Artist.ts'

export default interface Song {
  id: string
  title: string
  description: string
  isRecorded: boolean
  bpm?: number
  songsterrLink?: string
  youtubeLink?: string
  releaseDate?: string
  difficulty?: Difficulty
  imageUrl?: string
  lastTimePlayed?: string

  albumTrackNo?: number

  playlistTrackNo?: number
  playlistCreatedAt?: string

  album?: Album
  artist?: Artist
  guitarTuning?: GuitarTuning
  sections: SongSection[]
}

export interface SongSection {
  id: string
  name: string
  rehearsals: number
  songSectionType: SongSectionType
}

export interface SongSectionType {
  id: string
  name: string
}

export interface GuitarTuning {
  id: string
  name: string
}
