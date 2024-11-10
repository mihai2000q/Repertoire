import { z } from 'zod'

export interface SignUpForm {
  email: string
  password: string
  name: string
}

export const signUpValidation = z.object({
  email: z.string().email('Email is Invalid'),
  password: z
    .string()
    .trim()
    .min(8, 'Password should be at least 8 characters')
    .refine((val) => /[A-Z]/.test(val), 'Password should have at least 1 upper character')
    .refine((val) => /[a-z]/.test(val), 'Password should have at least 1 lower character')
    .refine((val) => /[0-9]/.test(val), 'Password should have at least 1 digit'),
  name: z.string().min(1, 'Name cannot be empty')
})
