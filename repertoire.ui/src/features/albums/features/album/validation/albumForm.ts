import { z } from 'zod/v4'
import { FileWithPath } from '@mantine/dropzone'
import {
  songsterrLinkValidator,
  youtubeLinkValidator
} from '../../../../../validation/custom/validators.ts'

export const addNewAlbumSongSchema = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
  songsterrLink: songsterrLinkValidator,
  youtubeLink: youtubeLinkValidator
})
export type AddNewAlbumSongForm = z.infer<typeof addNewAlbumSongSchema>

export const editAlbumHeaderSchema = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
  releaseDate: z.string().nullish(),
  image: z.string().or(z.object<FileWithPath>()).nullish(),
  artistId: z.string().optional()
})
export type EditAlbumHeaderForm = z.infer<typeof editAlbumHeaderSchema>
