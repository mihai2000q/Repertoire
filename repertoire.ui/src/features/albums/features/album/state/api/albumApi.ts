import { api } from '../../../../../../state/api.ts'
import HttpMessageResponse from '../../../../../../types/responses/HttpMessageResponse.ts'
import {
  AddSongsToAlbumRequest,
  MoveSongFromAlbumRequest,
  RemoveSongsFromAlbumRequest,
  UpdateAlbumRequest
} from '../../types/requests/AlbumRequests.ts'

const albumApi = api.injectEndpoints({
  endpoints: (build) => ({
    updateAlbum: build.mutation<HttpMessageResponse, UpdateAlbumRequest>({
      query: (body) => ({
        url: 'albums',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Albums']
    }),
    deleteImageFromAlbum: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `albums/images/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Albums']
    }),

    addSongsToAlbum: build.mutation<HttpMessageResponse, AddSongsToAlbumRequest>({
      query: (body) => ({
        url: `albums/add-songs`,
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Albums', 'Songs']
    }),
    moveSongFromAlbum: build.mutation<HttpMessageResponse, MoveSongFromAlbumRequest>({
      query: (body) => ({
        url: `albums/move-song`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Albums']
    }),
    removeSongsFromAlbum: build.mutation<HttpMessageResponse, RemoveSongsFromAlbumRequest>({
      query: (body) => ({
        url: `albums/remove-songs`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Albums', 'Songs']
    })
  })
})

export const {
  useUpdateAlbumMutation,
  useDeleteImageFromAlbumMutation,
  useAddSongsToAlbumMutation,
  useMoveSongFromAlbumMutation,
  useRemoveSongsFromAlbumMutation
} = albumApi
