import { z } from 'zod'

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
