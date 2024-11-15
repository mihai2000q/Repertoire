import { z } from 'zod'

const youtubeRegex =
  /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/(watch\?v=|embed\/|v\/|.+\?v=)?([^&=%?]{11})$/

export interface AddNewSongForm {
  title: string
  description: string

  albumId?: string
  albumTitle?: string
  artistId?: string
  artistName?: string

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
