import Album from './Album.ts'
import Song from './Song.ts'

export default interface Artist {
  id: string
  name: string
  isBand: boolean
  imageUrl?: string
  albums: Album[]
  songs: Song[]

  createdAt: string
  updatedAt: string
}
