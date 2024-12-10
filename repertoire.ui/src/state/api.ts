import { createApi } from '@reduxjs/toolkit/query/react'
import { queryWithRedirection } from './api.query'
import TokenResponse from '../types/responses/TokenResponse'
import User from '../types/models/User'
import { SignInRequest, SignUpRequest } from '../types/requests/AuthRequests.ts'

export const api = createApi({
  baseQuery: queryWithRedirection,
  reducerPath: 'api',
  tagTypes: ['Songs', 'Albums', 'Artists', 'Playlists', 'GuitarTunings', 'SongSectionTypes'],
  endpoints: (build) => {
    return {
      // Auth
      signUp: build.mutation<TokenResponse, SignUpRequest>({
        query: (body) => ({
          url: `auth/sign-up`,
          method: 'POST',
          body: body
        })
      }),
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
      })
    }
  }
})

export const { useSignUpMutation, useSignInMutation, useGetCurrentUserQuery } = api
