import { z } from 'zod'

export interface SignInForm {
  email: string
  password: string
}

export const signInValidation = z.object({
  email: z.string().email('Email is invalid'),
  password: z.string().trim().min(1, 'Password cannot be blank')
})
