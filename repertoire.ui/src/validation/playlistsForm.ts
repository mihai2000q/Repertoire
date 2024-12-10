import { z } from 'zod'

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
}

export const editPlaylistHeaderValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})
