import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Song, { GuitarTuning, Instrument, SongSectionType } from '../../types/models/Song.ts'
import {
  AddPartialSongRehearsalRequest,
  AddPerfectSongRehearsalRequest,
  CreateSongRequest,
  CreateSongSectionRequest,
  DeleteSongSectionRequest,
  GetSongsRequest,
  MoveSongSectionRequest,
  SaveImageToSongRequest,
  UpdateSongRequest,
  UpdateSongSectionRequest,
  UpdateSongSectionsOccurrencesRequest,
  UpdateSongSectionsPartialOccurrencesRequest
} from '../../types/requests/SongRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'

const songsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getSongs: build.query<WithTotalCountResponse<Song>, GetSongsRequest>({
      query: (arg) => ({
        url: `songs${createQueryParams(arg)}`
      }),
      providesTags: ['Songs']
    }),
    getSong: build.query<Song, string>({
      query: (arg) => `songs/${arg}`,
      providesTags: ['Songs'],
      transformResponse: (response: Song) => ({
        ...response,
        artist: response.artist
          ? {
              ...response.artist,
              bandMembers: response.artist.bandMembers === null ? [] : response.artist.bandMembers
            }
          : response.artist
      })
    }),
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
    addPartialSongRehearsal: build.mutation<HttpMessageResponse, AddPartialSongRehearsalRequest>({
      query: (body) => ({
        url: 'songs/partial-rehearsal',
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
      invalidatesTags: ['Songs', 'Albums']
    }),
    saveImageToSong: build.mutation<HttpMessageResponse, SaveImageToSongRequest>({
      query: (request) => ({
        url: 'songs/images',
        method: 'PUT',
        body: createFormData(request),
        formData: true
      }),
      invalidatesTags: ['Songs', 'Albums']
    }),
    deleteImageFromSong: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `songs/images/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Songs', 'Albums']
    }),
    deleteSong: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `songs/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Songs', 'Albums']
    }),

    // sections
    createSongSection: build.mutation<{ id: string }, CreateSongSectionRequest>({
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
    updateSongSectionsOccurrences: build.mutation<
      HttpMessageResponse,
      UpdateSongSectionsOccurrencesRequest
    >({
      query: (body) => ({
        url: 'songs/sections/occurrences',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    updateSongSectionsPartialOccurrences: build.mutation<
      HttpMessageResponse,
      UpdateSongSectionsPartialOccurrencesRequest
    >({
      query: (body) => ({
        url: 'songs/sections/partial-occurrences',
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
    deleteSongSection: build.mutation<HttpMessageResponse, DeleteSongSectionRequest>({
      query: (arg) => ({
        url: `songs/sections/${arg.id}/from/${arg.songId}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Songs']
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
  useCreateSongMutation,
  useAddPerfectSongRehearsalMutation,
  useAddPartialSongRehearsalMutation,
  useUpdateSongMutation,
  useSaveImageToSongMutation,
  useDeleteImageFromSongMutation,
  useDeleteSongMutation,
  useGetGuitarTuningsQuery,
  useGetInstrumentsQuery,
  useGetSongSectionTypesQuery,
  useCreateSongSectionMutation,
  useUpdateSongSectionMutation,
  useUpdateSongSectionsOccurrencesMutation,
  useUpdateSongSectionsPartialOccurrencesMutation,
  useMoveSongSectionMutation,
  useDeleteSongSectionMutation
} = songsApi
