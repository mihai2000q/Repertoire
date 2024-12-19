import { FileWithPath } from '@mantine/dropzone'
import { z } from 'zod'

export interface AccountForm {
  name: string
  profilePicture: string | FileWithPath | null
}

export const accountValidation = z.object({
  name: z.string().trim().min(1, 'Name cannot be blank')
})
