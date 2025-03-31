import { createApi } from '@reduxjs/toolkit/query/react'
import { SignInRequest } from '../../types/requests/AuthRequests.ts'
import { BaseQueryFn, FetchArgs, fetchBaseQuery, FetchBaseQueryError } from '@reduxjs/toolkit/query'
import { RootState } from '../store.ts'
import HttpErrorResponse from '../../types/responses/HttpErrorResponse.ts'
import { setErrorPath } from '../slice/globalSlice.ts'
import { toast } from 'react-toastify'

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

const errorCodeToPathname = new Map<number | string, string>([
  [403, '401'],
  [404, '404']
])

const authQueryWithRedirection: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError
> = async (args, api, extraOptions) => {
  const result = await authQueryWithAuthorization(args, api, extraOptions)
  const error = result?.error as HttpErrorResponse | undefined | null
  const errorStatus = error?.status ?? 0
  if (errorCodeToPathname.has(errorStatus)) {
    api.dispatch(setErrorPath(errorCodeToPathname.get(errorStatus)))
  } else if (errorStatus === 400) {
    toast.error(error.data.error)
  }
  return result
}

export const authApi = createApi({
  baseQuery: authQueryWithRedirection,
  reducerPath: 'authApi',
  endpoints: (build) => {
    return {
      signIn: build.mutation<string, SignInRequest>({
        query: (body) => ({
          url: `auth/sign-in`,
          method: 'PUT',
          body: body
        })
      }),
      getCentrifugeToken: build.query<string, void>({
        query: () => 'auth/centrifugo'
      })
    }
  }
})

export const { useSignInMutation, useLazyGetCentrifugeTokenQuery } = authApi
