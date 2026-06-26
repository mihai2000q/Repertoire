import { z } from 'zod/v4'

export const addNewPlaylistSchema = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
  description: z.string()
})
export type AddNewPlaylistForm = z.infer<typeof addNewPlaylistSchema>
