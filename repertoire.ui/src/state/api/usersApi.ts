import { api } from '../api.ts'
import { SignUpRequest } from '../../features/main/features/topbar/types/requests/UserRequests.ts'

const usersApi = api.injectEndpoints({
  endpoints: (build) => ({
    signUp: build.mutation<string, SignUpRequest>({
      query: (body) => ({
        url: `users/sign-up`,
        method: 'POST',
        body: body
      })
    })
  })
})

export const { useSignUpMutation } = usersApi
