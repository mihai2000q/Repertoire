import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import { SearchBase } from '../../types/models/Search.ts'
import { GetSearchRequest } from '../../types/requests/SearchRequests.ts'
import createQueryParams from '../../utils/createQueryParams.ts'

const searchApi = api.injectEndpoints({
  endpoints: (build) => ({
    getSearch: build.query<WithTotalCountResponse<SearchBase>, GetSearchRequest>({
      query: (arg) => `search${createQueryParams(arg)}`,
      providesTags: ['Artists', 'Albums', 'Songs', 'Playlists']
    })
  })
})

export const { useGetSearchQuery } = searchApi
