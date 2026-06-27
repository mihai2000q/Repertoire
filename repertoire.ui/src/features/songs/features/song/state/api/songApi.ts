import { UpdateSongRequest } from '../../types/requests/SongRequests.ts'
import HttpMessageResponse from '../../../../../../types/responses/HttpMessageResponse.ts'
import { api } from '../../../../../../state/api.ts'

const songApi = api.injectEndpoints({
  endpoints: (build) => ({
    updateSong: build.mutation<HttpMessageResponse, UpdateSongRequest>({
      query: (body) => ({
        url: 'songs',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    deleteImageFromSong: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `songs/images/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Songs']
    })
  })
})

export const { useUpdateSongMutation, useDeleteImageFromSongMutation } = songApi
