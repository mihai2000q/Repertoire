import { z } from 'zod'
import { FileWithPath } from '@mantine/dropzone'
import {
  songsterrLinkValidator,
  youtubeLinkValidator
} from './custom/validators.ts'

export interface AddNewSongForm {
  title: string
  description: string
  artistName?: string
  albumTitle?: string

  releaseDate?: string
  bpm?: number | string
  songsterrLink?: string
  youtubeLink?: string
}

export const addNewSongValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
  songsterrLink: songsterrLinkValidator,
  youtubeLink: youtubeLinkValidator
})

export interface EditSongHeaderForm {
  title: string
  releaseDate?: string
  image?: string | FileWithPath | null
  artistId?: string
  albumId?: string
}

export const editSongHeaderValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})

export interface EditSongLinksForm {
  songsterrLink?: string
  youtubeLink?: string
}

export const editSongLinksValidation = z.object({
  songsterrLink: songsterrLinkValidator,
  youtubeLink: youtubeLinkValidator
})

export interface EditSongSectionForm {
  name: string
  rehearsals: number | string
  confidence: number
  typeId: string
  bandMemberId?: string
  instrumentId?: string
}

export const editSongSectionValidation = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank')
})
