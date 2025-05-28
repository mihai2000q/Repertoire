import { z } from 'zod/v4'

export const youtubeRegex =
  /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/(watch\?v=|embed\/|v\/|.+\?v=)?([^&=%?]{11})$/

export const youtubeLinkValidator = z
  .preprocess(
    (input) => (input === '' ? null : input),
    z
      .url('This is not a valid URL')
      .trim()
      .regex(youtubeRegex, { error: 'This is not a valid Youtube URL' })
      .nullish()
  )
  .nullish()

export const songsterrLinkValidator = z
  .preprocess(
    (input) => (input === '' ? null : input),
    z
      .url('This is not a valid URL')
      .trim()
      .includes('songsterr.com', { error: 'This is not a valid Songsterr URL' })
      .nullish()
  )
  .nullish()
