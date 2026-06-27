import createQueryParams from '../../../../utils/createQueryParams.ts'
import { CreatePlaylistRequest } from '../../types/requests/PlaylistsRequests.ts'
import { PlaylistFiltersMetadata } from '../../types/PlaylistFilterMetadata.ts'
import { api } from '../../../../state/api.ts'

const playlistsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getPlaylistFiltersMetadata: build.query<PlaylistFiltersMetadata, { searchBy?: string[] }>({
      query: (arg) => `playlists/filters-metadata${createQueryParams(arg)}`,
      providesTags: ['Playlists', 'Songs']
    }),
    createPlaylist: build.mutation<{ id: string }, CreatePlaylistRequest>({
      query: (body) => ({
        url: 'playlists',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Playlists']
    })
  })
})

export const {
  useGetPlaylistFiltersMetadataQuery,
  useLazyGetPlaylistFiltersMetadataQuery,
  useCreatePlaylistMutation
} = playlistsApi
