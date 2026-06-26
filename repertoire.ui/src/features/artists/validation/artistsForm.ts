import { z } from 'zod/v4'

export const addNewArtistSchema = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank'),
  isBand: z.boolean()
})
export type AddNewArtistForm = z.infer<typeof addNewArtistSchema>
