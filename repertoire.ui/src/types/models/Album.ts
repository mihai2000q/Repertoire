import Song from './Song.ts'
import Artist from './Artist.ts'

export default interface Album {
  id: string
  title: string
  imageUrl?: string
  releaseDate?: string
  artist?: Artist
  songs: Song[]
  songsCount: number

  rehearsals: number
  confidence: number
  progress: number
  lastTimePlayed?: string

  createdAt: string
  updatedAt: string
}
