import { createApi } from '@reduxjs/toolkit/query/react'
import { queryWithRedirection } from './api.query'
import TokenResponse from '../types/responses/TokenResponse'
import HttpMessageResponse from '../types/responses/HttpMessageResponse.ts'
import User from '../types/models/User'
import { SignInRequest, SignUpRequest } from '../types/requests/AuthRequests.ts'
import {
  SaveProfilePictureToUserRequest,
  UpdateUserRequest
} from '../types/requests/UserRequests.ts'
import createFormData from '../utils/createFormData.ts'

export const api = createApi({
  baseQuery: queryWithRedirection,
  reducerPath: 'api',
  tagTypes: [
    'Songs',
    'Albums',
    'Artists',
    'Playlists',
    'GuitarTunings',
    'SongSectionTypes',
    'User'
  ],
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
        query: () => `users/current`,
        providesTags: ['User']
      }),
      updateUser: build.mutation<HttpMessageResponse, UpdateUserRequest>({
        query: (body) => ({
          url: 'users',
          method: 'PUT',
          body: body
        }),
        invalidatesTags: ['User']
      }),
      saveProfilePictureToUser: build.mutation<
        HttpMessageResponse,
        SaveProfilePictureToUserRequest
      >({
        query: (request) => ({
          url: 'users/pictures',
          method: 'PUT',
          body: createFormData(request),
          formData: true
        }),
        invalidatesTags: ['User']
      }),
      deleteProfilePictureFromUser: build.mutation<HttpMessageResponse, void>({
        query: () => ({
          url: 'users/pictures',
          method: 'DELETE'
        }),
        invalidatesTags: ['User']
      })
    }
  }
})

export const {
  useSignUpMutation,
  useSignInMutation,
  useGetCurrentUserQuery,
  useUpdateUserMutation,
  useSaveProfilePictureToUserMutation,
  useDeleteProfilePictureFromUserMutation
} = api
