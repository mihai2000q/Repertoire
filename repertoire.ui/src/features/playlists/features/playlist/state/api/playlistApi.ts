import HttpMessageResponse from '../../../../../../types/responses/HttpMessageResponse.ts'
import {
  MoveSongFromPlaylistRequest,
  RemoveSongsFromPlaylistRequest,
  ShufflePlaylistSongsRequest,
  UpdatePlaylistRequest
} from '../../types/requests/PlaylistRequests.ts'
import { api } from '../../../../../../state/api.ts'

const playlistApi = api.injectEndpoints({
  endpoints: (build) => ({
    updatePlaylist: build.mutation<HttpMessageResponse, UpdatePlaylistRequest>({
      query: (body) => ({
        url: 'playlists',
        method: 'PUT',
        body: body
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

    // songs
    shufflePlaylist: build.mutation<HttpMessageResponse, ShufflePlaylistSongsRequest>({
      query: (body) => ({
        url: 'playlists/songs/shuffle',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    moveSongFromPlaylist: build.mutation<HttpMessageResponse, MoveSongFromPlaylistRequest>({
      query: (body) => ({
        url: `playlists/songs/move`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    removeSongsFromPlaylist: build.mutation<HttpMessageResponse, RemoveSongsFromPlaylistRequest>({
      query: (body) => ({
        url: `playlists/songs/remove`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    })
  })
})

export const {
  useUpdatePlaylistMutation,
  useDeleteImageFromPlaylistMutation,
  useShufflePlaylistMutation,
  useMoveSongFromPlaylistMutation,
  useRemoveSongsFromPlaylistMutation
} = playlistApi
