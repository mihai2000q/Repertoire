import { z } from 'zod'
import { FileWithPath } from '@mantine/dropzone'

export interface AddNewArtistForm {
  name: string
}

export const addNewArtistValidation = z.object({
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

export interface EditArtistHeaderForm {
  name: string
  image: string | FileWithPath | null
}

export const editArtistHeaderValidation = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank')
})
