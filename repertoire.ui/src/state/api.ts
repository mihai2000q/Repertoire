import { createApi } from '@reduxjs/toolkit/query/react'
import { queryWithRedirection } from './api.query'
import TokenResponse from '../types/responses/TokenResponse'
import SignInRequest from '../types/requests/AuthRequests'
import User from '../types/models/User'

export const api = createApi({
  baseQuery: queryWithRedirection,
  reducerPath: 'api',
  tagTypes: ['Songs', 'Albums', 'Artists', 'GuitarTunings', 'SongSectionTypes'],
  endpoints: (build) => {
    return {
      // Auth
      signIn: build.mutation<TokenResponse, SignInRequest>({
        query: (body) => ({
          url: `auth/sign-in`,
          method: 'PUT',
          body: body
        })
      }),

      // Users
      getCurrentUser: build.query<User, void>({
        query: () => `users/current`
      }),
    }
  }
})

export const {
  useSignInMutation,
  useGetCurrentUserQuery
} = api
