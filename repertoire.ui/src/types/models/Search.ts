import SearchType from '../enums/SearchType.ts'

export interface SearchBase {
  id: string
  type: SearchType
}

export interface ArtistSearch extends SearchBase {
  name: string
  imageUrl?: string
}

export interface AlbumSearch extends SearchBase {
  title: string
  imageUrl?: string
  releaseDate?: string
  artist?: {
    id: string
    name: string
    imageUrl?: string
  }
}

export interface SongSearch extends SearchBase {
  title: string
  imageUrl?: string
  releaseDate?: string
  artist?: {
    id: string
    name: string
    imageUrl?: string
  }
  album?: {
    id: string
    title: string
    imageUrl?: string
  }
}

export interface PlaylistSearch extends SearchBase {
  title: string
  imageUrl?: string
}
