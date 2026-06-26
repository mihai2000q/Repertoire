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

  const drawersState = (api.getState() as RootState).drawers

  // null propagation for unit tests
  const isArtistDrawerOpenWithDeletedArtist =
    drawersState.artist?.artistId && args.url.includes(drawersState.artist.artistId)
  const isAlbumDrawerOpenWithDeletedAlbum =
    drawersState.album?.albumId && args.url.includes(drawersState.album.albumId)
  const isSongDrawerOpenWithDeletedSong =
    drawersState.song?.songId && args.url.includes(drawersState.song.songId)
  const isPlaylistDrawerOpenWithDeletedPlaylist =
    drawersState.playlist?.playlistId && args.url.includes(drawersState.playlist.playlistId)

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
