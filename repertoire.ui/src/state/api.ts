import { createApi } from '@reduxjs/toolkit/query/react'
import { queryWithRedirection } from './api.query'
import Song from '../types/models/Song'
import {
  CreateSongRequest,
  GetSongsRequest,
  UpdateSongRequest
} from '../types/requests/SongRequests'
import TokenResponse from '../types/responses/TokenResponse'
import SignInRequest from '../types/requests/AuthRequests'
import User from '../types/models/User'
import HttpMessageResponse from '../types/responses/HttpMessageResponse'

export const api = createApi({
  baseQuery: queryWithRedirection,
  reducerPath: 'api',
  tagTypes: ['Songs'],
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
      getCurrentUser: build.query<User, undefined>({
        query: () => `users/current`
      }),

      // Songs
      getSongs: build.query<{ songs: Song[], totalCount: number }, GetSongsRequest>({
        query: (arg) => ({
          url: 'songs',
          params: arg
        }),
        providesTags: ['Songs']
      }),
      getSong: build.query<Song, { id: string }>({
        query: (arg) => `songs${arg.id}`,
        providesTags: ['Songs']
      }),
      createSong: build.mutation<HttpMessageResponse, CreateSongRequest>({
        query: (body) => ({
          url: 'songs',
          method: 'POST',
          params: body
        }),
        invalidatesTags: ['Songs']
      }),
      updateSong: build.mutation<HttpMessageResponse, UpdateSongRequest>({
        query: (body) => ({
          url: 'songs',
          method: 'PUT',
          params: body
        }),
        invalidatesTags: ['Songs']
      }),
      deleteSong: build.mutation<HttpMessageResponse, { id: string }>({
        query: (arg) => ({
          url: `songs${arg.id}`,
          method: 'DELETE'
        }),
        invalidatesTags: ['Songs']
      })
    }
  }
})

export const {
  useSignInMutation,
  useGetCurrentUserQuery,
  useGetSongsQuery,
  useGetSongQuery,
  useCreateSongMutation,
  useUpdateSongMutation,
  useDeleteSongMutation
} = api