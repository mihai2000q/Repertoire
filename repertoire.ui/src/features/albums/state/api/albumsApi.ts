import { api } from '../../../../state/api.ts'
import { AlbumFiltersMetadata } from '../../types/AlbumFiltersMetadata.ts'
import createQueryParams from '../../../../utils/createQueryParams.ts'

const albumsApi = api.injectEndpoints({
  endpoints: (build) => {
    return {
      getAlbumFiltersMetadata: build.query<AlbumFiltersMetadata, { searchBy?: string[] }>({
        query: (arg) => `albums/filters-metadata${createQueryParams(arg)}`,
        providesTags: ['Albums', 'Artists', 'Songs'],
        transformResponse: (response: AlbumFiltersMetadata) => ({
          ...response,
          artistIds: response.artistIds ?? []
        })
      })
    }
  }
})

export const { useGetAlbumFiltersMetadataQuery, useLazyGetAlbumFiltersMetadataQuery } = albumsApi
