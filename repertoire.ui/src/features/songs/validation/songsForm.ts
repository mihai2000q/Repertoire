import { z } from 'zod/v4'
import { songsterrLinkValidator, youtubeLinkValidator } from '../../../validation/custom/validators.ts'

export const addNewSongSchema = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
  description: z.string(),
  artistName: z.string().optional(),
  albumTitle: z.string().optional(),

  releaseDate: z.string().optional(),
  bpm: z.number().or(z.string()).optional(),

  songsterrLink: songsterrLinkValidator,
  youtubeLink: youtubeLinkValidator
})
export type AddNewSongForm = z.infer<typeof addNewSongSchema>
