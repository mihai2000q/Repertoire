import { createApi } from '@reduxjs/toolkit/query/react'
import { queryWithRedirection } from './api.query'
import HttpMessageResponse from '../types/responses/HttpMessageResponse.ts'
import User from '../types/models/User'
import {
  SaveProfilePictureRequest,
  SignUpRequest,
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
    'Instruments',
    'BandMemberRoles',
    'SongSectionTypes',
    'User',
    'Search'
  ],
  endpoints: (build) => {
    return {
      getCurrentUser: build.query<User, void>({
        query: () => `users/current`,
        providesTags: ['User']
      }),
      signUp: build.mutation<string, SignUpRequest>({
        query: (body) => ({
          url: `sign-up`,
          method: 'POST',
          body: body
        })
      }),
      updateUser: build.mutation<HttpMessageResponse, UpdateUserRequest>({
        query: (body) => ({
          url: 'users',
          method: 'PUT',
          body: body
        }),
        invalidatesTags: ['User']
      }),
      deleteUser: build.mutation<HttpMessageResponse, void>({
        query: (body) => ({
          url: 'users',
          method: 'DELETE',
          body: body
        })
      }),

      saveProfilePicture: build.mutation<HttpMessageResponse, SaveProfilePictureRequest>({
        query: (request) => ({
          url: 'users/pictures',
          method: 'PUT',
          body: createFormData(request),
          formData: true
        }),
        invalidatesTags: ['User']
      }),
      deleteProfilePicture: build.mutation<HttpMessageResponse, void>({
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
  useGetCurrentUserQuery,
  useUpdateUserMutation,
  useDeleteUserMutation,
  useSaveProfilePictureMutation,
  useDeleteProfilePictureMutation
} = api
