import {z} from "zod";

export interface AddNewSongForm {
  title: string
}

export const addNewSongValidation = z.object({
  title: z.string().min(1)
})
