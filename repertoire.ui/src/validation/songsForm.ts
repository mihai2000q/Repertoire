import { z } from 'zod'

export interface AddNewSongForm {
  title: string
  description: string

  albumId?: string
  albumTitle?: string
  artistId?: string
  artistName?: string

  releaseDate?: Date
  bpm?: number
  songsterrLink?: string
  youtubeLink?: string
}

export const addNewSongValidation = z.object({
  title: z.string().trim().min(1, 'Title cannot be blank'),
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
      z.string().url('This is not a valid URL').nullish()
    )
    .nullish()
})
