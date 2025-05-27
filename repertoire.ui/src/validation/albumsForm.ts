import { z } from 'zod'
import { FileWithPath } from '@mantine/dropzone'
import {
  songsterrLinkValidator,
  youtubeLinkValidator
} from './custom/validators.ts'

export interface AddNewAlbumForm {
  title: string
  releaseDate?: string
  artistName?: string
}

export const addNewAlbumValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})

export interface AddNewAlbumSongForm {
  title: string
  songsterrLink?: string
  youtubeLink?: string
}

export const addNewAlbumSongValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
  songsterrLink: songsterrLinkValidator,
  youtubeLink: youtubeLinkValidator
})

export interface EditAlbumHeaderForm {
  title: string
  releaseDate?: string
  image?: string | FileWithPath | null
  artistId?: string
}

export const editAlbumHeaderValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})
