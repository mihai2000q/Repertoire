import {z} from "zod";

export interface AddNewAlbumForm {
  title: string
  releaseDate: string
}

export const addNewAlbumValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
})
