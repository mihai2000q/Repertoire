import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../types/models/Playlist.ts'
import {
  AddSongsToPlaylistRequest,
  CreatePlaylistRequest,
  GetPlaylistRequest,
  GetPlaylistsRequest,
  MoveSongFromPlaylistRequest,
  RemoveSongsFromPlaylistRequest,
  SaveImageToPlaylistRequest,
  UpdatePlaylistRequest
} from '../../types/requests/PlaylistRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'
import { PlaylistFiltersMetadata } from '../../types/models/FiltersMetadata.ts'

const playlistsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getPlaylists: build.query<WithTotalCountResponse<Playlist>, GetPlaylistsRequest>({
      query: (arg) => `playlists${createQueryParams(arg)}`,
      providesTags: ['Playlists', 'Songs']
    }),
    getPlaylist: build.query<Playlist, GetPlaylistRequest>({
      query: (arg) => `playlists/${arg.id}${createQueryParams({ ...arg, id: undefined })}`,
      providesTags: ['Playlists', 'Songs', 'Albums', 'Artists'],
      transformResponse: (response: Playlist) => ({
        ...response,
        songs: response.songs === null ? [] : response.songs
      })
    }),
    getPlaylistFiltersMetadata: build.query<PlaylistFiltersMetadata, { searchBy?: string[] }>({
      query: (arg) => `playlists/filters-metadata${createQueryParams(arg)}`,
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
      invalidatesTags: ['Songs', 'Playlists']
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
      invalidatesTags: ['Songs', 'Playlists']
    })
  })
})

export const {
  useGetPlaylistsQuery,
  useGetPlaylistQuery,
  useGetPlaylistFiltersMetadataQuery,
  useLazyGetPlaylistFiltersMetadataQuery,
  useCreatePlaylistMutation,
  useUpdatePlaylistMutation,
  useSaveImageToPlaylistMutation,
  useDeleteImageFromPlaylistMutation,
  useDeletePlaylistMutation,
  useAddSongsToPlaylistMutation,
  useMoveSongFromPlaylistMutation,
  useRemoveSongsFromPlaylistMutation
} = playlistsApi
