import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Song, {
  GuitarTuning,
  Instrument,
  SongArrangement,
  SongSectionType
} from '../../types/models/Song.ts'
import {
  AddPerfectSongRehearsalRequest,
  AddPerfectSongRehearsalsRequest,
  BulkDeleteSongSectionsRequest,
  BulkDeleteSongsRequest,
  BulkRehearsalsSongSectionsRequest,
  CreateSongArrangementRequest,
  CreateSongRequest,
  CreateSongSectionRequest,
  DeleteSongArrangementRequest,
  DeleteSongSectionRequest,
  GetSongArrangementsRequest,
  GetSongsRequest,
  MoveSongArrangementRequest,
  MoveSongSectionRequest,
  SaveImageToSongRequest,
  UpdateAllSongSectionsRequest,
  UpdateDefaultSongArrangementRequest,
  BulkUpdateSongArrangementsRequest,
  UpdateSongRequest,
  UpdateSongSectionRequest,
  UpdateSongSettingsRequest
} from '../../types/requests/SongRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'
import { SongFiltersMetadata } from '../../types/models/FiltersMetadata.ts'

const songsApi = api.injectEndpoints({
  endpoints: (build) => ({
    // Queries
    getSongs: build.query<WithTotalCountResponse<Song>, GetSongsRequest>({
      query: (arg) => `songs${createQueryParams(arg)}`,
      providesTags: ['Songs', 'Artists', 'Albums']
    }),
    getSong: build.query<Song, string>({
      query: (arg) => `songs/${arg}`,
      providesTags: ['Songs', 'Artists', 'Albums'],
      transformResponse: (response: Song) => ({
        ...response,
        artist: response.artist
          ? {
              ...response.artist,
              bandMembers: response.artist.bandMembers ?? []
            }
          : response.artist
      })
    }),
    getSongFiltersMetadata: build.query<SongFiltersMetadata, { searchBy?: string[] }>({
      query: (arg) => `songs/filters-metadata${createQueryParams(arg)}`,
      providesTags: ['Songs', 'Artists', 'Albums'],
      transformResponse: (response: SongFiltersMetadata) => ({
        ...response,
        artistIds: response.artistIds ?? [],
        albumIds: response.albumIds ?? []
      })
    }),

    // Infinite Queries
    getInfiniteSongs: build.infiniteQuery<
      WithTotalCountResponse<Song>,
      GetSongsRequest,
      { currentPage: number; pageSize: number }
    >({
      infiniteQueryOptions: {
        initialPageParam: {
          currentPage: 1,
          pageSize: 20
        },
        getNextPageParam: (lastPage, __, lastPageParam, ___, args) => {
          const pageSize = args.pageSize ?? lastPageParam.pageSize

          const totalSongs = lastPageParam.currentPage * pageSize
          const remainingSongs = lastPage?.totalCount - totalSongs

          if (remainingSongs <= 0) return undefined

          return {
            ...lastPageParam,
            currentPage: lastPageParam.currentPage + 1
          }
        }
      },
      query: ({ queryArg, pageParam }) => {
        const newQueryParams: GetSongsRequest = {
          ...queryArg,
          currentPage: pageParam.currentPage,
          pageSize: queryArg.pageSize ?? pageParam.pageSize
        }
        return `/songs${createQueryParams(newQueryParams)}`
      },
      providesTags: ['Songs', 'Albums', 'Artists']
    }),

    // Mutations
    createSong: build.mutation<{ id: string }, CreateSongRequest>({
      query: (body) => ({
        url: 'songs',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs', 'Artists', 'Albums']
    }),
    addPerfectSongRehearsal: build.mutation<HttpMessageResponse, AddPerfectSongRehearsalRequest>({
      query: (body) => ({
        url: 'songs/perfect-rehearsal',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    addPerfectSongRehearsals: build.mutation<HttpMessageResponse, AddPerfectSongRehearsalsRequest>({
      query: (body) => ({
        url: 'songs/perfect-rehearsals',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    updateSong: build.mutation<HttpMessageResponse, UpdateSongRequest>({
      query: (body) => ({
        url: 'songs',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    updateSongSettings: build.mutation<HttpMessageResponse, UpdateSongSettingsRequest>({
      query: (body) => ({
        url: 'songs/settings',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    bulkDeleteSongs: build.mutation<HttpMessageResponse, BulkDeleteSongsRequest>({
      query: (body) => ({
        url: `songs/bulk-delete`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    saveImageToSong: build.mutation<HttpMessageResponse, SaveImageToSongRequest>({
      query: (request) => ({
        url: 'songs/images',
        method: 'PUT',
        body: createFormData(request),
        formData: true
      }),
      invalidatesTags: ['Songs']
    }),
    deleteImageFromSong: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `songs/images/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Songs']
    }),
    deleteSong: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `songs/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Songs']
    }),

    // sections
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
    }),

    // arrangements
    getSongArrangements: build.query<SongArrangement[], GetSongArrangementsRequest>({
      query: (arg) => `songs/arrangements${createQueryParams(arg)}`,
      providesTags: ['SongArrangements', 'Songs']
    }),
    createSongArrangement: build.mutation<HttpMessageResponse, CreateSongArrangementRequest>({
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
    bulkUpdateSongArrangements: build.mutation<HttpMessageResponse, BulkUpdateSongArrangementsRequest>({
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
      invalidatesTags: ['SongArrangements']
    }),

    // sections - types
    getSongSectionTypes: build.query<SongSectionType[], void>({
      query: () => 'songs/sections/types',
      providesTags: ['SongSectionTypes']
    }),

    // guitar-tunings
    getGuitarTunings: build.query<GuitarTuning[], void>({
      query: () => 'songs/guitar-tunings',
      providesTags: ['GuitarTunings']
    }),

    // instruments
    getInstruments: build.query<Instrument[], void>({
      query: () => 'songs/instruments',
      providesTags: ['Instruments']
    })
  })
})

export const {
  useGetSongsQuery,
  useGetSongQuery,
  useGetSongFiltersMetadataQuery,
  useLazyGetSongFiltersMetadataQuery,
  useGetInfiniteSongsInfiniteQuery,
  useCreateSongMutation,
  useAddPerfectSongRehearsalMutation,
  useAddPerfectSongRehearsalsMutation,
  useUpdateSongMutation,
  useUpdateSongSettingsMutation,
  useBulkDeleteSongsMutation,
  useSaveImageToSongMutation,
  useDeleteImageFromSongMutation,
  useDeleteSongMutation,
  useGetGuitarTuningsQuery,
  useGetInstrumentsQuery,
  useGetSongSectionTypesQuery,
  useCreateSongSectionMutation,
  useBulkRehearsalsSongSectionsMutation,
  useUpdateSongSectionMutation,
  useUpdateAllSongSectionsMutation,
  useMoveSongSectionMutation,
  useBulkDeleteSongSectionsMutation,
  useDeleteSongSectionMutation,
  useGetSongArrangementsQuery,
  useCreateSongArrangementMutation,
  useBulkUpdateSongArrangementsMutation,
  useUpdateDefaultSongArrangementMutation,
  useMoveSongArrangementMutation,
  useDeleteSongArrangementMutation
} = songsApi
