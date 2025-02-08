import { FileWithPath } from '@mantine/dropzone'

export interface GetArtistsRequest {
  currentPage?: number
  pageSize?: number
  orderBy?: string[]
  searchBy?: string[]
}

export interface CreateArtistRequest {
  name: string
  isBand?: boolean
}

export interface UpdateArtistRequest {
  id: string
  name: string
  isBand?: boolean
}

export interface SaveImageToArtistRequest {
  image: FileWithPath
  id: string
}

export interface DeleteArtistRequest {
  id: string
  withAlbums?: boolean
  withSongs?: boolean
}

export interface AddAlbumsToArtistRequest {
  id: string
  albumIds: string[]
}

export interface RemoveAlbumsFromArtistRequest {
  id: string
  albumIds: string[]
}

export interface AddSongsToArtistRequest {
  id: string
  songIds: string[]
}

export interface RemoveSongsFromArtistRequest {
  id: string
  songIds: string[]
}

// Band Members

export interface CreateBandMemberRequest {
  name: string
  color?: string
  roleIds: string[]
  artistId: string
}

export interface UpdateBandMemberRequest {
  id: string
  name: string
  color?: string
  roleIds: string[]
}

export interface MoveBandMemberRequest {
  id: string
  overId: string
  artistId: string
}

export interface SaveImageToBandMemberRequest {
  image: FileWithPath
  id: string
}

export interface DeleteBandMemberRequest {
  id: string
  artistId: string
}
