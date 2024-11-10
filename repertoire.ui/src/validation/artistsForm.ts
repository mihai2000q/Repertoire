import { z } from 'zod'

export interface AddNewArtistForm {
  name: string
}

export const addNewArtistValidation = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank')
})
