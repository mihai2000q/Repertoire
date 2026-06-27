import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Song, {
  GuitarTuning,
  Instrument,
  SongArrangement,
  SongSectionType
} from '../../types/models/Song.ts'
import {
  AddCustomSongRehearsalRequest,
  AddCustomSongRehearsalsRequest,
  AddPerfectSongRehearsalRequest,
  AddPerfectSongRehearsalsRequest,
  BulkDeleteSongsRequest,
  CreateSongRequest,
  GetSongArrangementsRequest,
  GetSongsRequest,
  SaveImageToSongRequest,
  UpdateSongSettingsRequest
} from '../../types/requests/SongRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'

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
    addCustomSongRehearsal: build.mutation<HttpMessageResponse, AddCustomSongRehearsalRequest>({
      query: (body) => ({
        url: 'songs/custom-rehearsal',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    addCustomSongRehearsals: build.mutation<HttpMessageResponse, AddCustomSongRehearsalsRequest>({
      query: (body) => ({
        url: 'songs/custom-rehearsals',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
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

    getSongArrangements: build.query<SongArrangement[], GetSongArrangementsRequest>({
      query: (arg) => `songs/arrangements${createQueryParams(arg)}`,
      providesTags: ['SongArrangements', 'Songs']
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
  useGetInfiniteSongsInfiniteQuery,
  useCreateSongMutation,
  useAddCustomSongRehearsalMutation,
  useAddCustomSongRehearsalsMutation,
  useAddPerfectSongRehearsalMutation,
  useAddPerfectSongRehearsalsMutation,
  useUpdateSongSettingsMutation,
  useBulkDeleteSongsMutation,
  useSaveImageToSongMutation,
  useDeleteSongMutation,
  useGetSongArrangementsQuery,
  useGetGuitarTuningsQuery,
  useGetInstrumentsQuery,
  useGetSongSectionTypesQuery // move
} = songsApi
