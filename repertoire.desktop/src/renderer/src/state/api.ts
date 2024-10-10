import {
  BaseQueryFn,
  createApi,
  FetchArgs,
  fetchBaseQuery,
  FetchBaseQueryError
} from '@reduxjs/toolkit/query/react'
import { logout, setToken } from './authSlice'
import { RootState } from './store'
import { Mutex } from 'async-mutex'
import { setErrorPath } from './globalSlice'

const mutex = new Mutex()
const baseQuery = fetchBaseQuery({
  baseUrl: import.meta.env.VITE_BACKEND_URL,
  prepareHeaders: (headers, { getState }) => {
    const token = (getState() as RootState).auth.token
    if (token) {
      headers.set('authorization', `Bearer ${token}`)
    }
    return headers
  }
})

const queryWithRefreshToken: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
  args,
  api,
  extraOptions
) => {
  await mutex.waitForUnlock()
  let result = await baseQuery(args, api, extraOptions)
  if (result?.error?.status === 401 && !isAuthRequest(args)) {
    if (!mutex.isLocked()) {
      const release = await mutex.acquire()
      try {
        const authState = (api.getState() as RootState).auth
        const refreshResult = await baseQuery(
          {
            url: `auth/refresh`,
            method: 'PUT',
            body: { token: authState.token }
          },
          api,
          extraOptions
        )
        const data = refreshResult?.data as { token: string } | undefined
        if (data) {
          api.dispatch(setToken(data.token))
          result = await baseQuery(args, api, extraOptions)
        } else {
          api.dispatch(logout())
        }
      } finally {
        release()
      }
    } else {
      await mutex.waitForUnlock()
      result = await baseQuery(args, api, extraOptions)
    }
  }
  return result
}

const errorCodeToPathname = new Map<number | string, string>([
  [403, '401'],
  [404, '404']
])

const queryWithRedirection: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
  args,
  api,
  extraOptions
) => {
  const result = await queryWithRefreshToken(args, api, extraOptions)
  const error = result?.error as { status: number } | undefined | null
  const errorStatus = error?.status ?? 0
  if (errorCodeToPathname.has(errorStatus)) {
    api.dispatch(setErrorPath(errorCodeToPathname.get(errorStatus)))
  }
  return result
}

function isAuthRequest(args: string | FetchArgs): boolean {
  return (
    (typeof args === 'string' && args.includes('auth')) ||
    (typeof args === 'object' && args.url.includes('auth'))
  )
}

export const api = createApi({
  baseQuery: queryWithRedirection,
  reducerPath: 'api',
  tagTypes: [],
  endpoints: (build) => ({
    getCurrentUser: build.query({
      query: () => `users/current`
    })
  })
})

export const { useGetCurrentUserQuery } = api