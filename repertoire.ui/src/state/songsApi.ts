import { api } from './api'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse'
import Song from '../types/models/Song'
import {
  CreateSongRequest,
  GetSongsRequest,
  UpdateSongRequest
} from '../types/requests/SongRequests'
import HttpMessageResponse from '../types/responses/HttpMessageResponse'

const songsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getSongs: build.query<WithTotalCountResponse<Song>, GetSongsRequest>({
      query: (arg) => ({
        url: 'songs/',
        params: arg
      }),
      providesTags: ['Songs']
    }),
    getSong: build.query<Song, string>({
      query: (arg) => `songs/${arg}`,
      providesTags: ['Songs']
    }),
    createSong: build.mutation<HttpMessageResponse, CreateSongRequest>({
      query: (body) => ({
        url: 'songs/',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    updateSong: build.mutation<HttpMessageResponse, UpdateSongRequest>({
      query: (body) => ({
        url: 'songs/',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    deleteSong: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `songs/${arg}`,
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
  useDeleteSongMutation
} = songsApi
