import { z } from 'zod'
import { FileWithPath } from '@mantine/dropzone'

export interface AddNewPlaylistForm {
  title: string
  description: string
}

export const addNewPlaylistValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})

export interface EditPlaylistHeaderForm {
  title: string
  description: string
  image?: string | FileWithPath | null
}

export const editPlaylistHeaderValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})
