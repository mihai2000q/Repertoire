import { z } from 'zod/v4'
import { FileWithPath } from '@mantine/dropzone'
import {
  songsterrLinkValidator,
  youtubeLinkValidator
} from './custom/validators.ts'

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

export const editSongHeaderSchema = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
  releaseDate: z.string().nullish(),
  image: z.string().or(z.object<FileWithPath>()).nullish(),
  artistId: z.string().optional(),
  albumId: z.string().optional()
})
export type EditSongHeaderForm = z.infer<typeof editSongHeaderSchema>

export const editSongLinksSchema = z.object({
  songsterrLink: songsterrLinkValidator,
  youtubeLink: youtubeLinkValidator
})
export type EditSongLinksForm = z.infer<typeof editSongLinksSchema>

export const editSongSectionSchema = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank'),
  rehearsals: z.number().or(z.string()),
  confidence: z.number(),
  typeId: z.string(),
  bandMemberId: z.string().optional(),
  instrumentId: z.string().optional()
})
export type EditSongSectionForm = z.infer<typeof editSongSectionSchema>
