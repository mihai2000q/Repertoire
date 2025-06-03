import Song from './Song.ts'

export default interface Playlist {
  id: string
  title: string
  description: string
  imageUrl?: string
  songs: Song[]
  songsCount: number

  createdAt: string
  updatedAt: string
}
