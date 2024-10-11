import { createApi } from '@reduxjs/toolkit/query/react'
import { queryWithRedirection } from '@renderer/state/api.query'
import User from '@renderer/types/models/User'
import SignInRequest from '@renderer/types/requests/SignInRequest'
import TokenResponse from '@renderer/types/Token.response'

export const api = createApi({
  baseQuery: queryWithRedirection,
  reducerPath: 'api',
  tagTypes: [],
  endpoints: (build) => ({
    getCurrentUser: build.query({
      query: () => `users/current`
    })
  })
})

export const { useGetCurrentUserQuery } = api
