import { createApi } from '@reduxjs/toolkit/query/react'
import { queryWithRedirection } from './api.query'
import User from './types/models/User'
import SignInRequest from './types/requests/SignInRequest'
import TokenResponse from './types/Token.response'

export const api = createApi({
  baseQuery: queryWithRedirection,
  reducerPath: 'api',
  tagTypes: [],
  endpoints: (build) => ({
    // Auth
    signIn: build.mutation<TokenResponse, SignInRequest>({
      query: (body) => ({
        url: `auth/sign-in`,
        method: 'PUT',
        body: body
      })
    }),

    // Users
    getCurrentUser: build.query<User, undefined>({
      query: () => `users/current`
    })
  })
})

export const { useSignInMutation, useGetCurrentUserQuery } = api
