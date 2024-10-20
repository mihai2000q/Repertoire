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
import WithTotalCountResponse from "../types/responses/WithTotalCountResponse";

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
      getSongs: build.query<WithTotalCountResponse<Song>, GetSongsRequest>({
        query: (arg) => ({
          url: 'songs/',
          params: arg
        }),
        providesTags: ['Songs']
      }),
      getSong: build.query<Song, string>({
        query: (arg) => `songs/${arg}`,
        providesTags: ['Songs']
      }),
      createSong: build.mutation<HttpMessageResponse, CreateSongRequest>({
        query: (body) => ({
          url: 'songs/',
          method: 'POST',
          body: body
        }),
        invalidatesTags: ['Songs']
      }),
      updateSong: build.mutation<HttpMessageResponse, UpdateSongRequest>({
        query: (body) => ({
          url: 'songs/',
          method: 'PUT',
          body: body
        }),
        invalidatesTags: ['Songs']
      }),
      deleteSong: build.mutation<HttpMessageResponse, string>({
        query: (arg) => ({
          url: `songs/${arg}`,
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
