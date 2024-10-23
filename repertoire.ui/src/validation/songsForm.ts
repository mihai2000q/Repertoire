import {z} from "zod";

export interface AddNewSongForm {
  title: string
}

export const addNewSongValidation = z.object({
  title: z.string().trim().min(1, "Title cannot be blank"),
})
