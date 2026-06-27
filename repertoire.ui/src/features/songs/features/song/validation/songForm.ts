import { z } from 'zod/v4'
import { FileWithPath } from '@mantine/dropzone'
import {
  songsterrLinkValidator,
  youtubeLinkValidator
} from '../../../../../validation/custom/validators.ts'

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
