import { SignUpRequest } from '../../types/requests/SignUpRequests.ts'
import { api } from '../../../../../state/api.ts'

const signUpApi = api.injectEndpoints({
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

export const { useSignUpMutation } = signUpApi
