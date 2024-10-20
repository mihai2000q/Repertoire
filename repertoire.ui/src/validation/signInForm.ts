import { z } from 'zod'

export interface SignInForm {
  email: string
  password: string
}

export const signInValidation = z.object({
  email: z.string().min(1).email(),
  password: z.string().min(8)
})
