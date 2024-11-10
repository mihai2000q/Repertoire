import { api } from './api'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse'
import Album from '../types/models/Album.ts'
import { GetAlbumsRequest } from '../types/requests/AlbumRequests.ts'

const songsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getAlbums: build.query<WithTotalCountResponse<Album>, GetAlbumsRequest>({
      query: (arg) => ({
        url: 'albums',
        params: arg
      }),
      providesTags: ['Albums']
    })
  })
})

export const { useGetAlbumsQuery } = songsApi
