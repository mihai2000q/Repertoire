import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Album from '../../types/models/Album.ts'
import {
  AddPerfectRehearsalsToAlbumsRequest,
  AddSongsToAlbumRequest,
  CreateAlbumRequest,
  DeleteAlbumRequest,
  GetAlbumRequest,
  GetAlbumsRequest,
  MoveSongFromAlbumRequest,
  RemoveSongsFromAlbumRequest,
  SaveImageToAlbumRequest,
  UpdateAlbumRequest
} from '../../types/requests/AlbumRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'
import { AlbumFiltersMetadata } from '../../types/models/FiltersMetadata.ts'

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
    getAlbumFiltersMetadata: build.query<AlbumFiltersMetadata, { searchBy?: string[] }>({
      query: (arg) => `albums/filters-metadata${createQueryParams(arg)}`,
      providesTags: ['Albums', 'Artists', 'Songs'],
      transformResponse: (response: AlbumFiltersMetadata) => ({
        ...response,
        artistIds: response.artistIds ?? []
      })
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
    deleteImageFromAlbum: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `albums/images/${arg}`,
        method: 'DELETE'
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
    }),

    addSongsToAlbum: build.mutation<HttpMessageResponse, AddSongsToAlbumRequest>({
      query: (body) => ({
        url: `albums/add-songs`,
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Albums', 'Songs']
    }),
    moveSongFromAlbum: build.mutation<HttpMessageResponse, MoveSongFromAlbumRequest>({
      query: (body) => ({
        url: `albums/move-song`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Albums']
    }),
    removeSongsFromAlbum: build.mutation<HttpMessageResponse, RemoveSongsFromAlbumRequest>({
      query: (body) => ({
        url: `albums/remove-songs`,
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
  useGetAlbumFiltersMetadataQuery,
  useLazyGetAlbumFiltersMetadataQuery,
  useCreateAlbumMutation,
  useAddPerfectRehearsalsToAlbumsMutation,
  useUpdateAlbumMutation,
  useSaveImageToAlbumMutation,
  useDeleteImageFromAlbumMutation,
  useDeleteAlbumMutation,
  useAddSongsToAlbumMutation,
  useMoveSongFromAlbumMutation,
  useRemoveSongsFromAlbumMutation
} = albumsApi
