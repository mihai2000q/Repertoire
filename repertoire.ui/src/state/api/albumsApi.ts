import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Album from '../../types/models/Album.ts'
import {
  AddPerfectRehearsalsToAlbumsRequest,
  BulkDeleteAlbumsRequest,
  CreateAlbumRequest,
  DeleteAlbumRequest,
  GetAlbumRequest,
  GetAlbumsRequest,
  SaveImageToAlbumRequest
} from '../../types/requests/AlbumRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'

const albumsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getAlbums: build.query<WithTotalCountResponse<Album>, GetAlbumsRequest>({
      query: (arg) => `albums${createQueryParams(arg)}`,
      providesTags: ['Albums', 'Artists', 'Songs']
    }),
    getAlbum: build.query<Album, GetAlbumRequest>({
      query: (arg) => `albums/${arg.id}${createQueryParams({ ...arg, id: undefined })}`,
      providesTags: ['Albums', 'Artists', 'Songs']
    }),
    createAlbum: build.mutation<{ id: string }, CreateAlbumRequest>({
      query: (body) => ({
        url: 'albums',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Albums', 'Artists']
    }),
    addPerfectRehearsalsToAlbums: build.mutation<
      HttpMessageResponse,
      AddPerfectRehearsalsToAlbumsRequest
    >({
      query: (body) => ({
        url: 'albums/perfect-rehearsals',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Albums', 'Songs']
    }),
    bulkDeleteAlbums: build.mutation<HttpMessageResponse, BulkDeleteAlbumsRequest>({
      query: (body) => ({
        url: `albums/bulk-delete`,
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
    deleteAlbum: build.mutation<HttpMessageResponse, DeleteAlbumRequest>({
      query: (arg) => ({
        url: `albums/${arg.id}`,
        method: 'DELETE',
        params: { ...arg, id: undefined }
      }),
      invalidatesTags: ['Albums']
    })
  })
})

export const {
  useGetAlbumsQuery,
  useGetAlbumQuery,
  useCreateAlbumMutation,
  useAddPerfectRehearsalsToAlbumsMutation,
  useBulkDeleteAlbumsMutation,
  useSaveImageToAlbumMutation,
  useDeleteAlbumMutation
} = albumsApi
