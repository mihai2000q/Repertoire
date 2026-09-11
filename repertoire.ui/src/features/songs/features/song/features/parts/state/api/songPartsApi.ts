import { api } from '../../../../../../../../state/api.ts'
import HttpMessageResponse from '../../../../../../../../types/responses/HttpMessageResponse.ts'
import {
  BulkDeleteSongPartsRequest,
  BulkUpdateSongPartsRequest,
  CreateSongPartRequest,
  DeleteSongPartRequest,
  GetSongPartsRequest,
  MoveSongPartInSongRequest,
  UpdateAllSongPartsRequest,
  UpdateSongPartRequest
} from '../../types/requests/SongPartRequests.ts'
import { SongPart } from '../../../../../../../../types/models/Song.ts'
import createQueryParams from '../../../../../../../../utils/createQueryParams.ts'

const songPartsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getSongParts: build.query<SongPart[], GetSongPartsRequest>({
      query: (arg) => `songs/parts${createQueryParams(arg)}`,
      providesTags: ['Songs']
    }),
    createSongPart: build.mutation<HttpMessageResponse, CreateSongPartRequest>({
      query: (body) => ({
        url: 'songs/parts',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    updateSongPart: build.mutation<HttpMessageResponse, UpdateSongPartRequest>({
      query: (body) => ({
        url: 'songs/parts',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    updateAllSongParts: build.mutation<HttpMessageResponse, UpdateAllSongPartsRequest>({
      query: (body) => ({
        url: 'songs/parts/all',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    moveSongPartInSong: build.mutation<HttpMessageResponse, MoveSongPartInSongRequest>({
      query: (body) => ({
        url: 'songs/parts/move-in-song',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    bulkUpdateSongParts: build.mutation<HttpMessageResponse, BulkUpdateSongPartsRequest>({
      query: (body) => ({
        url: 'songs/parts/bulk-update',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    bulkDeleteSongParts: build.mutation<HttpMessageResponse, BulkDeleteSongPartsRequest>({
      query: (body) => ({
        url: 'songs/parts/bulk-delete',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    deleteSongPart: build.mutation<HttpMessageResponse, DeleteSongPartRequest>({
      query: (arg) => ({
        url: `songs/parts/${arg.id}/from/${arg.songId}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Songs']
    })
  })
})

export const {
  useGetSongPartsQuery,
  useCreateSongPartMutation,
  useUpdateSongPartMutation,
  useUpdateAllSongPartsMutation,
  useMoveSongPartInSongMutation,
  useBulkUpdateSongPartsMutation,
  useBulkDeleteSongPartsMutation,
  useDeleteSongPartMutation
} = songPartsApi
