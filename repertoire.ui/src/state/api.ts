import { createApi } from '@reduxjs/toolkit/query/react'
import { queryWithRedirection } from './api.query'
import TokenResponse from '../types/responses/TokenResponse'
import HttpMessageResponse from '../types/responses/HttpMessageResponse.ts'
import User from '../types/models/User'
import { SignInRequest, SignUpRequest } from '../types/requests/AuthRequests.ts'
import { SaveProfilePictureRequest, UpdateUserRequest } from '../types/requests/UserRequests.ts'
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
      getCentrifugeToken: build.query<string, void>({
        query: () => 'auth/centrifugo',
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
      deleteUser: build.mutation<HttpMessageResponse, void>({
        query: (body) => ({
          url: 'users',
          method: 'DELETE',
          body: body
        }),
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
  useSignInMutation,
  useGetCentrifugeTokenQuery,
  useLazyGetCentrifugeTokenQuery,
  useGetCurrentUserQuery,
  useUpdateUserMutation,
  useDeleteUserMutation,
  useSaveProfilePictureMutation,
  useDeleteProfilePictureMutation
} = api
