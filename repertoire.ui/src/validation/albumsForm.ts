import { z } from 'zod'
import { FileWithPath } from '@mantine/dropzone'

export interface AddNewAlbumForm {
  title: string
  releaseDate?: Date
  artistName?: string
}

export const addNewAlbumValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})

export interface AddNewAlbumSongForm {
  title: string
}

export const addNewAlbumSongValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})

export interface EditAlbumHeaderForm {
  title: string
  releaseDate: Date
  image: string | FileWithPath | null
}

export const editAlbumHeaderValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})
