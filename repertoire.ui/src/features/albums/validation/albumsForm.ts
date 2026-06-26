import { z } from 'zod/v4'

export const addNewAlbumSchema = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank').default(''),
  releaseDate: z.string().nullish(),
  artistName: z.string().optional()
})
export type AddNewAlbumForm = z.infer<typeof addNewAlbumSchema>
