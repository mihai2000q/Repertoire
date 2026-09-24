import { api } from '../../../../../../state/api.ts'
import HttpMessageResponse from '../../../../../../types/responses/HttpMessageResponse.ts'
import {
  BulkDeleteSongSectionsRequest,
  CreateSongSectionRequest,
  DeleteSongSectionRequest,
  GetSongSectionsRequest,
  MoveSongSectionRequest,
  UpdateSongSectionRequest
} from '../../types/requests/SongSectionRequests.ts'
import { SongSection } from '../../../../../../types/models/Song.ts'
import createQueryParams from '../../../../../../utils/createQueryParams.ts'

const songSectionsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getSongSections: build.query<SongSection[], GetSongSectionsRequest>({
      query: (arg) => `songs/sections${createQueryParams(arg)}`,
      providesTags: ['Songs']
    }),
    createSongSection: build.mutation<HttpMessageResponse, CreateSongSectionRequest>({
      query: (body) => ({
        url: 'songs/sections',
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
        url:
          `songs/sections/${arg.id}/from/${arg.songId}` +
          `${createQueryParams({ ...arg, id: undefined, songId: undefined })}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Songs']
    })
  })
})

export const {
  useGetSongSectionsQuery,
  useCreateSongSectionMutation,
  useUpdateSongSectionMutation,
  useMoveSongSectionMutation,
  useBulkDeleteSongSectionsMutation,
  useDeleteSongSectionMutation
} = songSectionsApi
