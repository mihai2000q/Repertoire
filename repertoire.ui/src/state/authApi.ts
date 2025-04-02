import { createApi } from '@reduxjs/toolkit/query/react'
import { SignInRequest } from '../types/requests/AuthRequests.ts'
import { authQueryWithRedirection } from './authApi.query.ts'

export const authApi = createApi({
  baseQuery: authQueryWithRedirection,
  reducerPath: 'authApi',
  endpoints: (build) => {
    return {
      signIn: build.mutation<string, SignInRequest>({
        query: (body) => ({
          url: `sign-in`,
          method: 'PUT',
          body: body
        })
      }),
      getCentrifugeToken: build.query<string, void>({
        query: () => 'centrifugo/token'
      })
    }
  }
})

export const { useSignInMutation, useLazyGetCentrifugeTokenQuery } = authApi
