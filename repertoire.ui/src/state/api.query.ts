import { Mutex } from 'async-mutex'
import { RootState } from './store'
import { setToken, signOut } from './slice/authSlice.ts'
import {
  deleteAlbumDrawer,
  deleteArtistDrawer,
  deleteSongDrawer,
  setErrorPath
} from './slice/globalSlice.ts'
import {
  BaseQueryApi,
  BaseQueryFn,
  FetchArgs,
  fetchBaseQuery,
  FetchBaseQueryError
} from '@reduxjs/toolkit/query'
import { toast } from 'react-toastify'
import HttpErrorResponse from '../types/responses/HttpErrorResponse.ts'
import { DeleteArtistRequest } from '../types/requests/ArtistRequests.ts'
import { DeleteAlbumRequest } from '../types/requests/AlbumRequests.ts'

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

const refreshQuery = fetchBaseQuery({
  baseUrl: import.meta.env.VITE_AUTH_URL
})

const mutex = new Mutex()
const queryWithRefresh: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
  args,
  api,
  extraOptions
) => {
  await mutex.waitForUnlock()
  let result = await queryWithAuthorization(args, api, extraOptions)
  if (result?.error?.status === 401 && !(typeof args === 'object' && !args.url.includes('users/sign-up'))) {
    if (!mutex.isLocked()) {
      const release = await mutex.acquire()
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
  if (!error) cleanDrawers(args, api)
  const errorStatus = error?.status ?? 0
  if (errorCodeToPathname.has(errorStatus)) {
    api.dispatch(setErrorPath(errorCodeToPathname.get(errorStatus)))
  } else if (errorStatus === 400) {
    toast.error(error.data.error)
  }
  return result
}

function cleanDrawers(args: string | FetchArgs, api: BaseQueryApi) {
  if (typeof args !== 'object' || args.method !== 'DELETE') return

  const globalState = (api.getState() as RootState).global

  // null propagation for unit tests
  const isArtistDrawerOpenWithDeletedArtist =
    globalState.artistDrawer?.artistId && args.url.includes(globalState.artistDrawer.artistId)
  const isAlbumDrawerOpenWithDeletedAlbum =
    globalState.albumDrawer?.albumId && args.url.includes(globalState.albumDrawer.albumId)
  const isSongDrawerOpenWithDeletedSong =
    globalState.songDrawer?.songId && args.url.includes(globalState.songDrawer.songId)

  if (api.endpoint === 'deleteArtist') {
    if (isArtistDrawerOpenWithDeletedArtist) api.dispatch(deleteArtistDrawer())
    if ((args.params as DeleteArtistRequest)?.withAlbums === true) api.dispatch(deleteAlbumDrawer())
    if ((args.params as DeleteArtistRequest)?.withSongs === true) api.dispatch(deleteSongDrawer())
  } else if (api.endpoint === 'deleteAlbum') {
    if (isAlbumDrawerOpenWithDeletedAlbum) api.dispatch(deleteAlbumDrawer())
    if ((args.params as DeleteAlbumRequest)?.withSongs === true) api.dispatch(deleteSongDrawer())
  } else if (api.endpoint === 'deleteSong' && isSongDrawerOpenWithDeletedSong) {
    api.dispatch(deleteSongDrawer())
  }
}
