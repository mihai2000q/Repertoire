import { z } from 'zod'
import { FileWithPath } from '@mantine/dropzone'

export interface AddNewArtistForm {
  name: string
  isBand: boolean
}

export const addNewArtistValidation = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank')
})

export interface EditArtistHeaderForm {
  name: string
  image: string | FileWithPath | null
  isBand: boolean
}

export const editArtistHeaderValidation = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank')
})

export interface AddNewBandMemberForm {
  name: string
  image: string | FileWithPath | null
}

export const addNewBandMemberValidation = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank')
})

export interface AddNewArtistAlbumForm {
  title: string
  releaseDate?: string
}

export interface EditBandMemberForm {
  name: string
  color?: string
  image: string | FileWithPath | null
  roleIds: string[]
}

export const editBandMemberValidation = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank')
})

export interface AddNewArtistAlbumForm {
  title: string
}

export const addNewArtistAlbumValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})

export interface AddNewArtistSongForm {
  title: string
}

export const addNewArtistSongValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})
