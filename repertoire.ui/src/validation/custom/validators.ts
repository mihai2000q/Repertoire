import { z } from 'zod'

export const youtubeRegex =
  /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/(watch\?v=|embed\/|v\/|.+\?v=)?([^&=%?]{11})$/

export const youtubeLinkValidator = z
  .preprocess(
    (input) => (input === '' ? null : input),
    z
      .string()
      .url('This is not a valid URL')
      .regex(youtubeRegex, { message: 'This is not a valid Youtube URL' })
      .nullish()
  )
  .nullish()

export const songsterrLinkValidator = z
  .preprocess(
    (input) => (input === '' ? null : input),
    z
      .string()
      .includes('songsterr.com', { message: 'This is not a valid Songsterr URL' })
      .url('This is not a valid URL')
      .nullish()
  )
  .nullish()
