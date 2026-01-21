import { api } from '../api.ts'
import User from '../../types/models/User.ts'
import {
  SaveProfilePictureRequest,
  SignUpRequest,
  UpdateUserRequest
} from '../../types/requests/UserRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'

const usersApi = api.injectEndpoints({
  endpoints: (build) => ({
    getCurrentUser: build.query<User, void>({
      query: () => `users/current`,
      providesTags: ['User']
    }),
    signUp: build.mutation<string, SignUpRequest>({
      query: (body) => ({
        url: `users/sign-up`,
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
  })
})

export const {
  useSignUpMutation,
  useGetCurrentUserQuery,
  useUpdateUserMutation,
  useDeleteUserMutation,
  useSaveProfilePictureMutation,
  useDeleteProfilePictureMutation
} = usersApi
