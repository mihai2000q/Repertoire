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
  UpdateAllSongSectionsRequest,
  UpdateSongRequest,
  UpdateSongSectionRequest,
  UpdateSongSectionsOccurrencesRequest,
  UpdateSongSectionsPartialOccurrencesRequest,
  UpdateSongSettingsRequest
} from '../../types/requests/SongRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'
import { SongFiltersMetadata } from '../../types/models/FiltersMetadata.ts'

const songsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getSongs: build.query<WithTotalCountResponse<Song>, GetSongsRequest>({
      query: (arg) => ({
        url: `songs${createQueryParams(arg)}`
      }),
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
              bandMembers: response.artist.bandMembers === null ? [] : response.artist.bandMembers
            }
          : response.artist
      })
    }),
    getSongFiltersMetadata: build.query<SongFiltersMetadata, void>({
      query: () => 'songs/filters-metadata',
      providesTags: ['Songs', 'Artists', 'Albums'],
      transformResponse: (response: SongFiltersMetadata) => ({
        ...response,
        artistIds: response.artistIds === null ? [] : response.artistIds,
        albumIds: response.albumIds === null ? [] : response.albumIds
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
  useLazyGetSongFiltersMetadataQuery,
  useCreateSongMutation,
  useAddPerfectSongRehearsalMutation,
  useAddPartialSongRehearsalMutation,
  useUpdateSongMutation,
  useUpdateSongSettingsMutation,
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
  useUpdateAllSongSectionsMutation,
  useMoveSongSectionMutation,
  useDeleteSongSectionMutation
} = songsApi
