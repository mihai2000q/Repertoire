import { api } from '../../../../state/api.ts'
import createQueryParams from '../../../../utils/createQueryParams.ts'
import { ArtistFiltersMetadata } from '../../types/ArtistFilterMetadata.ts'
import { CreateArtistRequest } from '../../types/requests/ArtistsRequests.ts'

const artistsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getArtistFiltersMetadata: build.query<ArtistFiltersMetadata, { searchBy?: string[] }>({
      query: (arg) => `artists/filters-metadata${createQueryParams(arg)}`,
      providesTags: ['Albums', 'Artists', 'Songs']
    }),
    createArtist: build.mutation<{ id: string }, CreateArtistRequest>({
      query: (body) => ({
        url: 'artists',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Artists']
    })
  })
})

export const {
  useGetArtistFiltersMetadataQuery,
  useLazyGetArtistFiltersMetadataQuery,
  useCreateArtistMutation
} = artistsApi
