import Album from './Album.ts'
import Song from './Song.ts'

export default interface Artist {
  id: string
  name: string
  isBand: boolean
  imageUrl?: string
  albums: Album[]
  songs: Song[]
  bandMembers: BandMember[]

  songsCount: number
  rehearsals: number
  confidence: number
  progress: number
  lastTimePlayed?: string

  createdAt: string
  updatedAt: string
}

export interface BandMember {
  id: string
  name: string
  imageUrl?: string
  color?: string
  roles: BandMemberRole[]
}

export interface BandMemberRole {
  id: string
  name: string
}
