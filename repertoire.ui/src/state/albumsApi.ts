import { api } from './api'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse'
import Album from '../types/models/Album.ts'
import {
  CreateAlbumRequest,
  GetAlbumsRequest,
  SaveImageToAlbumRequest,
  UpdateAlbumRequest
} from '../types/requests/AlbumRequests.ts'
import HttpMessageResponse from '../types/responses/HttpMessageResponse.ts'
import createFormData from '../utils/createFormData.ts'

const albumsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getAlbums: build.query<WithTotalCountResponse<Album>, GetAlbumsRequest>({
      query: (arg) => ({
        url: 'albums',
        params: arg
      }),
      providesTags: ['Albums']
    }),
    getAlbum: build.query<Album, string>({
      query: (arg) => `albums/${arg}`,
      providesTags: ['Albums']
    }),
    createAlbum: build.mutation<{ id: string }, CreateAlbumRequest>({
      query: (body) => ({
        url: 'albums',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Albums']
    }),
    updateAlbum: build.mutation<HttpMessageResponse, UpdateAlbumRequest>({
      query: (body) => ({
        url: 'albums',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Albums']
    }),
    saveImageToAlbum: build.mutation<HttpMessageResponse, SaveImageToAlbumRequest>({
      query: (request) => ({
        url: 'albums/images',
        method: 'PUT',
        body: createFormData(request),
        formData: true
      }),
      invalidatesTags: ['Albums']
    }),
    deleteAlbum: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `albums/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Albums']
    })
  })
})

export const {
  useGetAlbumsQuery,
  useGetAlbumQuery,
  useCreateAlbumMutation,
  useUpdateAlbumMutation,
  useSaveImageToAlbumMutation,
  useDeleteAlbumMutation
} = albumsApi
