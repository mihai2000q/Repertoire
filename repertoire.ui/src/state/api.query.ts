import { Mutex } from 'async-mutex'
import { RootState } from './store.ts'
import { setToken, signOut } from './slice/authSlice.ts'
import { setErrorPath } from './slice/globalSlice.ts'
import {
  BaseQueryApi,
  BaseQueryFn,
  FetchArgs,
  fetchBaseQuery,
  FetchBaseQueryError
} from '@reduxjs/toolkit/query'
import { toast } from 'react-toastify'
import HttpErrorResponse from '../types/responses/HttpErrorResponse.ts'

type Arguments = string | FetchArgs
type ApiQuery = BaseQueryFn<Arguments, unknown, FetchBaseQueryError>

export const queryWithAuthorization = (baseUrl: string) =>
  fetchBaseQuery({
    baseUrl: baseUrl,
    prepareHeaders: (headers, { getState }) => {
      const token = (getState() as RootState).auth.token
      if (token) {
        headers.set('authorization', `Bearer ${token}`)
      }
      return headers
    }
  })

const refreshQuery = fetchBaseQuery({
  baseUrl: import.meta.env.VITE_AUTH_URL
})

export const apiMutex = new Mutex()
export const queryWithRefresh: (
  baseUrl: string,
  shouldRefresh?: (args: Arguments) => boolean
) => ApiQuery =
  (baseUrl: string, shouldRefresh: (args: Arguments) => boolean) =>
  async (args, api, extraOptions) => {
    await apiMutex.waitForUnlock()
    let result = await queryWithAuthorization(baseUrl)(args, api, extraOptions)
    if (result?.error?.status === 401 && shouldRefresh && shouldRefresh(args)) {
      if (!apiMutex.isLocked()) {
        const release = await apiMutex.acquire()
        try {
          const authState = (api.getState() as RootState).auth
          const refreshResult = await refreshQuery(
            {
              url: `refresh`,
              method: 'PUT',
              body: { token: authState.token }
            },
            api,
            extraOptions
          )
          const newToken = refreshResult?.data as string | undefined
          if (newToken) {
            api.dispatch(setToken(newToken))
            result = await queryWithAuthorization(baseUrl)(args, api, extraOptions)
          } else {
            api.dispatch(signOut())
          }
        } finally {
          release()
        }
      } else {
        await apiMutex.waitForUnlock()
        result = await queryWithAuthorization(baseUrl)(args, api, extraOptions)
      }
    }
    return result
  }

const errorCodeToPathname = new Map<number, string>([
  [403, '401'],
  [404, '404']
])

export const queryWithRedirection: (
  apiQuery: ApiQuery,
  onError?: (args: Arguments, api: BaseQueryApi) => void
) => ApiQuery =
  (apiQuery: ApiQuery, onError: (args: Arguments, api: BaseQueryApi) => void) =>
  async (args, api, extraOptions) => {
    const result = await apiQuery(args, api, extraOptions)
    const error = result?.error as HttpErrorResponse | undefined | null

    if (!error && onError) onError(args, api)

    const errorStatus = error?.status ?? 0
    if (errorCodeToPathname.has(errorStatus)) {
      api.dispatch(setErrorPath(errorCodeToPathname.get(errorStatus)))
    } else if (error && errorStatus !== 401) {
      toast.error(error.data?.error)
    }
    return result
  }
