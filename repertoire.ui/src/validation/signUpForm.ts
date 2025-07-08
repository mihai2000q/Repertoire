import { z } from 'zod/v4'

export const signUpSchema = z.object({
  email: z.email('Email is invalid'),
  password: z
    .string()
    .refine((val) => /[A-Z]/.test(val), 'Password must have at least 1 upper character')
    .refine((val) => /[a-z]/.test(val), 'Password must have at least 1 lower character')
    .refine((val) => /[0-9]/.test(val), 'Password must have at least 1 digit')
    .refine((val) => val.length >= 8, 'Password must have at least 8 characters')
    .refine((val) => val.length > 0, 'Password cannot be blank'),
  name: z.string().trim().min(1, 'Name cannot be blank')
})
export type SignUpForm = z.infer<typeof signUpSchema>
