import { createApi } from '@reduxjs/toolkit/query/react'
import { queryWithRedirection, queryWithRefresh } from './api.query'
import { BaseQueryApi, FetchArgs } from '@reduxjs/toolkit/query'
import { RootState } from './store.ts'
import {
  deleteAlbumDrawer,
  deleteArtistDrawer,
  deletePlaylistDrawer,
  deleteSongDrawer
} from './slice/drawersSlice.ts'
import { DeleteArtistRequest } from '../types/requests/ArtistRequests.ts'
import { DeleteAlbumRequest } from '../types/requests/AlbumRequests.ts'

function cleanDrawers(args: string | FetchArgs, api: BaseQueryApi) {
  if (typeof args !== 'object' || args.method !== 'DELETE') return

  const globalState = (api.getState() as RootState).drawers

  // null propagation for unit tests
  const isArtistDrawerOpenWithDeletedArtist =
    globalState.artist?.artistId && args.url.includes(globalState.artist.artistId)
  const isAlbumDrawerOpenWithDeletedAlbum =
    globalState.album?.albumId && args.url.includes(globalState.album.albumId)
  const isSongDrawerOpenWithDeletedSong =
    globalState.song?.songId && args.url.includes(globalState.song.songId)
  const isPlaylistDrawerOpenWithDeletedPlaylist =
    globalState.playlist?.playlistId &&
    args.url.includes(globalState.playlist.playlistId)

  if (api.endpoint === 'deleteArtist') {
    if (isArtistDrawerOpenWithDeletedArtist) api.dispatch(deleteArtistDrawer())
    if ((args.params as DeleteArtistRequest)?.withAlbums === true) api.dispatch(deleteAlbumDrawer())
    if ((args.params as DeleteArtistRequest)?.withSongs === true) api.dispatch(deleteSongDrawer())
  } else if (api.endpoint === 'deleteAlbum') {
    if (isAlbumDrawerOpenWithDeletedAlbum) api.dispatch(deleteAlbumDrawer())
    if ((args.params as DeleteAlbumRequest)?.withSongs === true) api.dispatch(deleteSongDrawer())
  } else if (api.endpoint === 'deleteSong' && isSongDrawerOpenWithDeletedSong) {
    api.dispatch(deleteSongDrawer())
  } else if (api.endpoint === 'deletePlaylist' && isPlaylistDrawerOpenWithDeletedPlaylist) {
    api.dispatch(deletePlaylistDrawer())
  }
}

const query = queryWithRedirection(
  queryWithRefresh(
    import.meta.env.VITE_BACKEND_URL,
    (args) => !(typeof args === 'object' && args.url.includes('users/sign-up'))
  ),
  (args, a) => cleanDrawers(args, a)
)

export const api = createApi({
  baseQuery: query,
  reducerPath: 'api',
  tagTypes: [
    'Albums',

    'Artists',
    'BandMemberRoles',

    'Playlists',

    'Search',

    'Songs',
    'SongArrangements',
    'SongSectionTypes',
    'GuitarTunings',
    'Instruments',

    'User'
  ],
  keepUnusedDataFor: 5 * 60,
  refetchOnReconnect: true,
  endpoints: () => ({})
})
