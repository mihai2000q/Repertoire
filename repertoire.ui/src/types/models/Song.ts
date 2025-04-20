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

  rehearsals: number
  confidence: number
  progress: number

  albumTrackNo?: number

  playlistTrackNo?: number
  playlistCreatedAt?: string

  settings: SongSettings
  album?: Album
  artist?: Artist
  guitarTuning?: GuitarTuning
  sections: SongSection[]

  solosCount: number
  riffsCount: number

  createdAt: string
  updatedAt: string
}

export interface SongSettings {
  id: string
  defaultBandMember?: BandMember
  defaultInstrument?: Instrument
}

export interface SongSection {
  id: string
  name: string
  rehearsals: number
  confidence: number
  progress: number
  occurrences: number
  partialOccurrences: number
  songSectionType: SongSectionType
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
