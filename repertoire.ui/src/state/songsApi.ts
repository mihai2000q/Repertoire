import { api } from './api'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse'
import Song, { GuitarTuning, SongSectionType } from '../types/models/Song'
import {
  CreateSongRequest,
  CreateSongSectionRequest,
  DeleteSongSectionRequest,
  GetSongsRequest,
  MoveSongSectionRequest,
  SaveImageToSongRequest,
  UpdateSongRequest,
  UpdateSongSectionRequest
} from '../types/requests/SongRequests'
import HttpMessageResponse from '../types/responses/HttpMessageResponse'
import createFormData from '../utils/createFormData.ts'
import createQueryParams from '../utils/createQueryParams.ts'

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
      providesTags: ['Songs']
    }),
    createSong: build.mutation<{ id: string }, CreateSongRequest>({
      query: (body) => ({
        url: 'songs',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs', 'Artists', 'Albums']
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
    deleteSong: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `songs/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Songs', 'Albums']
    }),

    getGuitarTunings: build.query<GuitarTuning[], void>({
      query: () => 'songs/guitar-tunings',
      providesTags: ['GuitarTunings']
    }),

    getSongSectionTypes: build.query<SongSectionType[], void>({
      query: () => 'songs/sections/types',
      providesTags: ['SongSectionTypes']
    }),
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
  useGetSongSectionTypesQuery,
  useCreateSongSectionMutation,
  useUpdateSongSectionMutation,
  useMoveSongSectionMutation,
  useDeleteSongSectionMutation
} = songsApi
