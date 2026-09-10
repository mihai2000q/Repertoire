import Difficulty from '../enums/Difficulty.ts'
import Album from './Album.ts'
import Artist, { BandMember } from './Artist.ts'

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
  defaultArrangementId?: string

  rehearsals: number
  confidence: number
  progress: number

  albumTrackNo?: number

  playlistSongId?: string
  playlistTrackNo?: number
  playlistCreatedAt?: string

  settings: SongSettings
  album?: Album
  artist?: Artist
  guitarTuning?: GuitarTuning
  parts: SongPart[]
  arrangements: SongArrangement[]

  solosCount?: number
  riffsCount?: number

  createdAt: string
  updatedAt: string
}

export interface SongSettings {
  id: string
  defaultBandMember?: BandMember
  defaultInstrument?: Instrument
}

export interface SongPart {
  id: string
  name: string
  rehearsals: number
  confidence: number
  progress: number
  sectionIds: string[] // TODO: change to sections
  bandMember?: BandMember
  instrument?: Instrument
}

export interface SongSectionType {
  id: string
  name: string
}

export interface GuitarTuning {
  id: string
  name: string
}

export interface Instrument {
  id: string
  name: string
}

export interface SongArrangement {
  id: string
  name: string
  partOccurrences: SongPartOccurrences[]
  songId: string
}

export interface SongPartOccurrences {
  section: SongPart
  occurrences: number | string // string, for front end purposes and transformations
}
