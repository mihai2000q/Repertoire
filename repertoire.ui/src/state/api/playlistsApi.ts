import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../types/models/Playlist.ts'
import {
  AddSongsToPlaylistRequest,
  CreatePlaylistRequest,
  GetPlaylistsRequest,
  MoveSongFromPlaylistRequest,
  RemoveSongsFromPlaylistRequest,
  SaveImageToPlaylistRequest,
  UpdatePlaylistRequest
} from '../../types/requests/PlaylistRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'

const playlistsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getPlaylists: build.query<WithTotalCountResponse<Playlist>, GetPlaylistsRequest>({
      query: (arg) => ({
        url: `playlists${createQueryParams(arg)}`
      }),
      providesTags: ['Playlists']
    }),
    getPlaylist: build.query<Playlist, string>({
      query: (arg) => `playlists/${arg}`,
      providesTags: ['Playlists', 'Songs']
    }),
    createPlaylist: build.mutation<{ id: string }, CreatePlaylistRequest>({
      query: (body) => ({
        url: 'playlists',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Playlists']
    }),
    updatePlaylist: build.mutation<HttpMessageResponse, UpdatePlaylistRequest>({
      query: (body) => ({
        url: 'playlists',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Playlists']
    }),
    saveImageToPlaylist: build.mutation<HttpMessageResponse, SaveImageToPlaylistRequest>({
      query: (request) => ({
        url: 'playlists/images',
        method: 'PUT',
        body: createFormData(request),
        formData: true
      }),
      invalidatesTags: ['Playlists']
    }),
    deleteImageFromPlaylist: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `playlists/images/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Playlists']
    }),
    deletePlaylist: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `playlists/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Playlists']
    }),

    addSongsToPlaylist: build.mutation<HttpMessageResponse, AddSongsToPlaylistRequest>({
      query: (body) => ({
        url: `playlists/add-songs`,
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Playlists']
    }),
    moveSongFromPlaylist: build.mutation<HttpMessageResponse, MoveSongFromPlaylistRequest>({
      query: (body) => ({
        url: `playlists/move-song`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Playlists']
    }),
    removeSongsFromPlaylist: build.mutation<HttpMessageResponse, RemoveSongsFromPlaylistRequest>({
      query: (body) => ({
        url: `playlists/remove-songs`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Playlists']
    })
  })
})

export const {
  useGetPlaylistsQuery,
  useGetPlaylistQuery,
  useCreatePlaylistMutation,
  useUpdatePlaylistMutation,
  useSaveImageToPlaylistMutation,
  useDeleteImageFromPlaylistMutation,
  useDeletePlaylistMutation,
  useAddSongsToPlaylistMutation,
  useMoveSongFromPlaylistMutation,
  useRemoveSongsFromPlaylistMutation
} = playlistsApi
