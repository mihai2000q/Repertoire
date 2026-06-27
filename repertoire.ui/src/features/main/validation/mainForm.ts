import { FileWithPath } from '@mantine/dropzone'
import { z } from 'zod/v4'

export const accountSchema = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank'),
  profilePicture: z.string().or(z.object<FileWithPath>()).nullish()
})
export type AccountForm = z.infer<typeof accountSchema>

export const deleteAccountSchema = z.object({
  password: z.string().trim().min(1, 'Password cannot be blank')
})
export type DeleteAccountForm = z.infer<typeof deleteAccountSchema>
