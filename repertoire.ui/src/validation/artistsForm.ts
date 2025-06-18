import { z } from 'zod/v4'
import { FileWithPath } from '@mantine/dropzone'
import { songsterrLinkValidator, youtubeLinkValidator } from './custom/validators.ts'

export const addNewArtistSchema = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank'),
  isBand: z.boolean()
})
export type AddNewArtistForm = z.infer<typeof addNewArtistSchema>

export const editArtistHeaderSchema = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank'),
  image: z.string().or(z.object<FileWithPath>()).nullish(),
  isBand: z.boolean()
})
export type EditArtistHeaderForm = z.infer<typeof editArtistHeaderSchema>

export const addNewBandMemberSchema = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank'),
  image: z.string().or(z.object<FileWithPath>()).nullish(),
  roleIds: z.array(z.string()).min(1, 'Select at least one role')
})
export type AddNewBandMemberForm = z.infer<typeof addNewBandMemberSchema>

export const editBandMemberSchema = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank'),
  image: z.string().or(z.object<FileWithPath>()).nullish(),
  color: z.string().nullish(),
  roleIds: z.array(z.string()).min(1, 'Select at least one role')
})
export type EditBandMemberForm = z.infer<typeof editBandMemberSchema>

export const addNewArtistAlbumSchema = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})
export type AddNewArtistAlbumForm = z.infer<typeof addNewArtistAlbumSchema>

export const addNewArtistSongSchema = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
  songsterrLink: songsterrLinkValidator,
  youtubeLink: youtubeLinkValidator
})
export type AddNewArtistSongForm = z.infer<typeof addNewArtistSongSchema>
