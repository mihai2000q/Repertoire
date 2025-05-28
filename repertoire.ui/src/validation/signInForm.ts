import { z } from 'zod/v4'

export const signInSchema = z.object({
  email: z.email('Email is invalid'),
  password: z.string().min(1, 'Password cannot be blank')
})
export type SignInForm = z.infer<typeof signInSchema>
