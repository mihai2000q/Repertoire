import { createApi } from '@reduxjs/toolkit/query/react'
import { SignInRequest } from '../types/requests/AuthRequests.ts'
import { queryWithRedirection, queryWithRefresh } from './api.query.ts'

const query = queryWithRedirection(
  queryWithRefresh(
    import.meta.env.VITE_AUTH_URL,
    (args) =>
      (typeof args === 'string' && args.includes('centrifugo')) ||
      (typeof args === 'object' && args.url.includes('centrifugo'))
  )
)

export const authApi = createApi({
  baseQuery: query,
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
