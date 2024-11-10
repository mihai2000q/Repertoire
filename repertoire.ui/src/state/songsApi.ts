import { api } from './api'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse'
import Song, { GuitarTuning, SongSectionType } from '../types/models/Song'
import {
  CreateSongRequest,
  GetSongsRequest,
  SaveImageToSongRequest,
  UpdateSongRequest
} from '../types/requests/SongRequests'
import HttpMessageResponse from '../types/responses/HttpMessageResponse'
import createFormData from '../utils/createFormData.ts'

const songsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getSongs: build.query<WithTotalCountResponse<Song>, GetSongsRequest>({
      query: (arg) => ({
        url: 'songs',
        params: arg
      }),
      providesTags: ['Songs']
    }),
    getSong: build.query<Song, string>({
      query: (arg) => `songs/${arg}`,
      providesTags: ['Songs']
    }),
    createSong: build.mutation<{ id: string }, CreateSongRequest>({
      query: (body) => ({
        url: 'songs',
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
    saveImageToSong: build.mutation<HttpMessageResponse, SaveImageToSongRequest>({
      query: (request) => ({
        url: 'songs/images',
        method: 'PUT',
        body: createFormData(request),
        formData: true
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

    getGuitarTunings: build.query<GuitarTuning[], void>({
      query: () => 'songs/guitar-tunings',
      providesTags: ['GuitarTunings']
    }),
    getSongSectionTypes: build.query<SongSectionType[], void>({
      query: () => 'songs/sections/types',
      providesTags: ['SongSectionTypes']
    })
  })
})

export const {
  useGetSongsQuery,
  useGetSongQuery,
  useCreateSongMutation,
  useUpdateSongMutation,
  useSaveImageToSongMutation,
  useDeleteSongMutation,
  useGetGuitarTuningsQuery,
  useGetSongSectionTypesQuery
} = songsApi
