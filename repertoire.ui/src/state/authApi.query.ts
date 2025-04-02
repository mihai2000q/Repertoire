import { RootState } from './store'
import { setToken, signOut } from './slice/authSlice.ts'
import { setErrorPath } from './slice/globalSlice.ts'
import { BaseQueryFn, FetchArgs, fetchBaseQuery, FetchBaseQueryError } from '@reduxjs/toolkit/query'
import { toast } from 'react-toastify'
import HttpErrorResponse from '../types/responses/HttpErrorResponse.ts'
import { apiMutex } from './api.query.ts'

const authQueryWithAuthorization = fetchBaseQuery({
  baseUrl: import.meta.env.VITE_AUTH_URL,
  prepareHeaders: (headers, { getState }) => {
    const token = (getState() as RootState).auth.token
    if (token) {
      headers.set('authorization', `Bearer ${token}`)
    }
    return headers
  }
})

const authQueryWithRefresh: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
  args,
  api,
  extraOptions
) => {
  await apiMutex.waitForUnlock()
  let result = await authQueryWithAuthorization(args, api, extraOptions)
  if (result?.error?.status === 401 && isCentrifugoRequest(args)) {
    if (!apiMutex.isLocked()) {
      const release = await apiMutex.acquire()
      try {
        const authState = (api.getState() as RootState).auth
        const refreshResult = await authQueryWithAuthorization(
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
          result = await authQueryWithAuthorization(args, api, extraOptions)
        } else {
          api.dispatch(signOut())
        }
      } finally {
        release()
      }
    } else {
      await apiMutex.waitForUnlock()
      result = await authQueryWithAuthorization(args, api, extraOptions)
    }
  }
  return result
}

const errorCodeToPathname = new Map<number | string, string>([
  [403, '401'],
  [404, '404']
])

export const authQueryWithRedirection: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError
> = async (args, api, extraOptions) => {
  const result = await authQueryWithRefresh(args, api, extraOptions)
  const error = result?.error as HttpErrorResponse | undefined | null
  const errorStatus = error?.status ?? 0
  if (errorCodeToPathname.has(errorStatus)) {
    api.dispatch(setErrorPath(errorCodeToPathname.get(errorStatus)))
  } else if (errorStatus === 400) {
    toast.error(error.data.error)
  }
  return result
}

function isCentrifugoRequest(args: string | FetchArgs): boolean {
  return (
    (typeof args === 'string' && args.includes('centrifugo')) ||
    (typeof args === 'object' && args.url.includes('centrifugo'))
  )
}
