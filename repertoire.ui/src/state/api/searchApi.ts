import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import { SearchBase } from '../../types/models/Search.ts'
import { GetSearchRequest } from '../../types/requests/SearchRequests.ts'
import createQueryParams from '../../utils/createQueryParams.ts'

const searchApi = api.injectEndpoints({
  endpoints: (build) => ({
    getSearch: build.query<WithTotalCountResponse<SearchBase>, GetSearchRequest>({
      query: (arg) => `search${createQueryParams(arg)}`,
      providesTags: ['Search'],
      transformResponse: (response: WithTotalCountResponse<SearchBase>) => ({
        ...response,
        models: response.models ?? []
      })
    }),

    getInfiniteSearch: build.infiniteQuery<
      WithTotalCountResponse<SearchBase>,
      GetSearchRequest,
      { currentPage: number; pageSize: number }
    >({
      infiniteQueryOptions: {
        initialPageParam: {
          currentPage: 1,
          pageSize: 20
        },
        getNextPageParam: (lastPage, __, lastPageParam, ___, args) => {
          const pageSize = args.pageSize ?? lastPageParam.pageSize

          const totalItems = lastPageParam.currentPage * pageSize
          const remainingItems = lastPage?.totalCount - totalItems

          if (remainingItems <= 0) return undefined

          return {
            ...lastPageParam,
            currentPage: lastPageParam.currentPage + 1
          }
        }
      },
      query: ({ queryArg, pageParam }) => {
        const newQueryParams: GetSearchRequest = {
          ...queryArg,
          currentPage: pageParam.currentPage,
          pageSize: queryArg.pageSize ?? pageParam.pageSize
        }
        return `search${createQueryParams(newQueryParams)}`
      },
      providesTags: ['Search']
    })
  })
})

export const { useGetSearchQuery, useGetInfiniteSearchInfiniteQuery } = searchApi
