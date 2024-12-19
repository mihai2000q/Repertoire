import { z } from 'zod'
import {FileWithPath} from "@mantine/dropzone";

const youtubeRegex =
  /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/(watch\?v=|embed\/|v\/|.+\?v=)?([^&=%?]{11})$/

export interface AddNewSongForm {
  title: string
  description: string
  artistName?: string
  albumTitle?: string

  releaseDate?: Date
  bpm?: number | string
  songsterrLink?: string
  youtubeLink?: string
}

export const addNewSongValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
  songsterrLink: z
    .string()
    .includes('songsterr.com', { message: 'This is not a valid Songsterr URL' })
    .url('This is not a valid URL')
    .nullish(),
  youtubeLink: z
    .string()
    .url('This is not a valid URL')
    .regex(youtubeRegex, { message: 'This is not a valid Youtube URL' })
    .nullish()
})

export interface EditSongHeaderForm {
  title: string
  releaseDate: Date
  image: string | FileWithPath | null
}

export const editSongHeaderValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank')
})

export interface EditSongLinksForm {
  songsterrLink?: string
  youtubeLink?: string
}

export const editSongLinksValidation = z.object({
  songsterrLink: z
    .preprocess(
      (input) => (input === '' ? null : input),
      z
        .string()
        .includes('songsterr.com', { message: 'This is not a valid Songsterr URL' })
        .url('This is not a valid URL')
        .nullish()
    )
    .nullish(),
  youtubeLink: z
    .preprocess(
      (input) => (input === '' ? null : input),
      z
        .string()
        .url('This is not a valid URL')
        .regex(youtubeRegex, { message: 'This is not a valid Youtube URL' })
        .nullish()
    )
    .nullish()
})

export interface EditSongSectionForm {
  name: string
  rehearsals: number | string
  confidence: number
  typeId: string
}

export const editSongSectionValidation = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank')
})
