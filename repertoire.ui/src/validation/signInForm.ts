import { z } from 'zod'

export interface SignInForm {
  email: string
  password: string
}

export const signInValidation = z.object({
  email: z.string().email(),
  password: z.string().min(8)
})
