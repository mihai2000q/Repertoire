import { Mutex } from 'async-mutex'
import { RootState } from './store'
import { setToken, signOut } from './slice/authSlice.ts'
import { setErrorPath } from './slice/globalSlice.ts'
import { BaseQueryFn, FetchArgs, fetchBaseQuery, FetchBaseQueryError } from '@reduxjs/toolkit/query'
import { toast } from 'react-toastify'
import HttpErrorResponse from '../types/responses/HttpErrorResponse.ts'

const queryWithAuthorization = fetchBaseQuery({
  baseUrl: import.meta.env.VITE_BACKEND_URL,
  prepareHeaders: (headers, { getState }) => {
    const token = (getState() as RootState).auth.token
    if (token) {
      headers.set('authorization', `Bearer ${token}`)
    }
    return headers
  }
})

const mutex = new Mutex()
const queryWithRefresh: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
  args,
  api,
  extraOptions
) => {
  await mutex.waitForUnlock()
  let result = await queryWithAuthorization(args, api, extraOptions)
  if (result?.error?.status === 401 && !isAuthRequest(args)) {
    if (!mutex.isLocked()) {
      const release = await mutex.acquire()
      try {
        const authState = (api.getState() as RootState).auth
        const refreshResult = await queryWithAuthorization(
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
          result = await queryWithAuthorization(args, api, extraOptions)
        } else {
          api.dispatch(signOut())
        }
      } finally {
        release()
      }
    } else {
      await mutex.waitForUnlock()
      result = await queryWithAuthorization(args, api, extraOptions)
    }
  }
  return result
}

const errorCodeToPathname = new Map<number | string, string>([
  [403, '401'],
  [404, '404']
])

export const queryWithRedirection: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError
> = async (args, api, extraOptions) => {
  const result = await queryWithRefresh(args, api, extraOptions)
  const error = result?.error as HttpErrorResponse | undefined | null
  const errorStatus = error?.status ?? 0
  if (errorCodeToPathname.has(errorStatus)) {
    api.dispatch(setErrorPath(errorCodeToPathname.get(errorStatus)))
  } else if (errorStatus === 400) {
    toast.error(error.data.error)
  }
  return result
}

function isAuthRequest(args: string | FetchArgs): boolean {
  return (
    (typeof args === 'string' && args.includes('auth')) ||
    (typeof args === 'object' && args.url.includes('auth'))
  )
}
