import { api } from '../../../../../../../../state/api.ts'
import { SongArrangement } from '../../../../../../../../types/models/Song.ts'
import HttpMessageResponse from '../../../../../../../../types/responses/HttpMessageResponse.ts'
import createQueryParams from '../../../../../../../../utils/createQueryParams.ts'
import {
  BulkUpdateSongArrangementsRequest,
  CreateSongArrangementRequest,
  DeleteSongArrangementRequest,
  GetSongArrangementsRequest,
  MoveSongArrangementRequest,
  UpdateDefaultSongArrangementRequest
} from '../../types/requests/SongArrangementsRequests.ts'

const songArrangementsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getSongArrangements: build.query<SongArrangement[], GetSongArrangementsRequest>({
      query: (arg) => `songs/arrangements${createQueryParams(arg)}`,
      providesTags: ['SongArrangements', 'Songs']
    }),
    createSongArrangement: build.mutation<{ id: string }, CreateSongArrangementRequest>({
      query: (body) => ({
        url: 'songs/arrangements',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['SongArrangements']
    }),
    moveSongArrangement: build.mutation<HttpMessageResponse, MoveSongArrangementRequest>({
      query: (body) => ({
        url: 'songs/arrangements/move',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['SongArrangements']
    }),
    bulkUpdateSongArrangements: build.mutation<
      HttpMessageResponse,
      BulkUpdateSongArrangementsRequest
    >({
      query: (body) => ({
        url: 'songs/arrangements/bulk',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['SongArrangements']
    }),
    updateDefaultSongArrangement: build.mutation<
      HttpMessageResponse,
      UpdateDefaultSongArrangementRequest
    >({
      query: (body) => ({
        url: 'songs/arrangements/default',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['SongArrangements', 'Songs']
    }),
    deleteSongArrangement: build.mutation<HttpMessageResponse, DeleteSongArrangementRequest>({
      query: (arg) => ({
        url: `songs/arrangements/${arg.id}/from/${arg.songId}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['SongArrangements', 'Songs']
    })
  })
})

export const {
  useGetSongArrangementsQuery,
  useCreateSongArrangementMutation,
  useBulkUpdateSongArrangementsMutation,
  useUpdateDefaultSongArrangementMutation,
  useMoveSongArrangementMutation,
  useDeleteSongArrangementMutation
} = songArrangementsApi
