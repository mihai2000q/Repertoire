import { z } from 'zod/v4'
import { FileWithPath } from '@mantine/dropzone'

export const editPlaylistHeaderSchema = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
  description: z.string(),
  image: z.string().or(z.object<FileWithPath>()).nullish()
})
export type EditPlaylistHeaderForm = z.infer<typeof editPlaylistHeaderSchema>
