import { FileWithPath } from '@mantine/dropzone'
import { z } from 'zod'

export interface AccountForm {
  name: string
  profilePicture?: string | FileWithPath | null
}

export const accountValidation = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank')
})

export interface DeleteAccountForm {
  password: string
}

export const deleteAccountValidation = z.object({
  password: z.string().trim().min(1, 'Password cannot be blank')
})
