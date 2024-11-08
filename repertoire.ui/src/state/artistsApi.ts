import { api } from './api'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse'
import Artist from '../types/models/Artist.ts'
import { GetArtistsRequest } from '../types/requests/ArtistRequests.ts'

const songsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getArtists: build.query<WithTotalCountResponse<Artist>, GetArtistsRequest>({
      query: (arg) => ({
        url: 'artists',
        params: arg
      }),
      providesTags: ['Artists']
    })
  })
})

export const { useGetArtistsQuery } = songsApi
