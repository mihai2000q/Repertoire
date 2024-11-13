import { api } from './api'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse'
import Album from '../types/models/Album.ts'
import {
  AddSongsToAlbumRequest,
  CreateAlbumRequest,
  GetAlbumsRequest,
  RemoveSongsFromAlbumRequest,
  SaveImageToAlbumRequest,
  UpdateAlbumRequest
} from '../types/requests/AlbumRequests.ts'
import HttpMessageResponse from '../types/responses/HttpMessageResponse.ts'
import createFormData from '../utils/createFormData.ts'
import createQueryParams from '../utils/createQueryParams.ts'

const albumsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getAlbums: build.query<WithTotalCountResponse<Album>, GetAlbumsRequest>({
      query: (arg) => ({
        url: `albums${createQueryParams(arg)}`
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
    }),

    addSongsToAlbum: build.mutation<HttpMessageResponse, AddSongsToAlbumRequest>({
      query: (body) => ({
        url: `artists/add-songs`,
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Albums', 'Songs']
    }),
    removeSongsFromAlbum: build.mutation<HttpMessageResponse, RemoveSongsFromAlbumRequest>({
      query: (body) => ({
        url: `artists/remove-songs`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Albums', 'Songs']
    })
  })
})

export const {
  useGetAlbumsQuery,
  useGetAlbumQuery,
  useCreateAlbumMutation,
  useUpdateAlbumMutation,
  useSaveImageToAlbumMutation,
  useDeleteAlbumMutation,
  useAddSongsToAlbumMutation,
  useRemoveSongsFromAlbumMutation
} = albumsApi
