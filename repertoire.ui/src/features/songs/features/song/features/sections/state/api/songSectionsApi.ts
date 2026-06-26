import { api } from '../../../../../../../../state/api.ts'
import HttpMessageResponse from '../../../../../../../../types/responses/HttpMessageResponse.ts'
import {
  BulkDeleteSongSectionsRequest,
  BulkRehearsalsSongSectionsRequest,
  CreateSongSectionRequest,
  DeleteSongSectionRequest,
  MoveSongSectionRequest,
  UpdateAllSongSectionsRequest,
  UpdateSongSectionRequest
} from '../../types/requests/SongSectionRequests.ts'

const songSectionsApi = api.injectEndpoints({
  endpoints: (build) => ({
    createSongSection: build.mutation<HttpMessageResponse, CreateSongSectionRequest>({
      query: (body) => ({
        url: 'songs/sections',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    bulkRehearsalsSongSections: build.mutation<
      HttpMessageResponse,
      BulkRehearsalsSongSectionsRequest
    >({
      query: (body) => ({
        url: 'songs/sections/bulk-rehearsals',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    updateSongSection: build.mutation<HttpMessageResponse, UpdateSongSectionRequest>({
      query: (body) => ({
        url: 'songs/sections',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    updateAllSongSections: build.mutation<HttpMessageResponse, UpdateAllSongSectionsRequest>({
      query: (body) => ({
        url: 'songs/sections/all',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    moveSongSection: build.mutation<HttpMessageResponse, MoveSongSectionRequest>({
      query: (body) => ({
        url: 'songs/sections/move',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    bulkDeleteSongSections: build.mutation<HttpMessageResponse, BulkDeleteSongSectionsRequest>({
      query: (body) => ({
        url: 'songs/sections/bulk-delete',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    deleteSongSection: build.mutation<HttpMessageResponse, DeleteSongSectionRequest>({
      query: (arg) => ({
        url: `songs/sections/${arg.id}/from/${arg.songId}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Songs']
    })
  })
})

export const {
  useCreateSongSectionMutation,
  useBulkRehearsalsSongSectionsMutation,
  useUpdateSongSectionMutation,
  useUpdateAllSongSectionsMutation,
  useMoveSongSectionMutation,
  useBulkDeleteSongSectionsMutation,
  useDeleteSongSectionMutation
} = songSectionsApi
