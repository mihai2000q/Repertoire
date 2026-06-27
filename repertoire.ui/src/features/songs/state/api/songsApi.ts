import { api } from '../../../../state/api.ts'
import { SongFiltersMetadata } from '../../types/SongFilterMetadata.ts'
import createQueryParams from '../../../../utils/createQueryParams.ts'

const songsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getSongFiltersMetadata: build.query<SongFiltersMetadata, { searchBy?: string[] }>({
      query: (arg) => `songs/filters-metadata${createQueryParams(arg)}`,
      providesTags: ['Songs', 'Artists', 'Albums'],
      transformResponse: (response: SongFiltersMetadata) => ({
        ...response,
        artistIds: response.artistIds ?? [],
        albumIds: response.albumIds ?? []
      })
    })
  })
})

export const { useGetSongFiltersMetadataQuery, useLazyGetSongFiltersMetadataQuery } = songsApi
