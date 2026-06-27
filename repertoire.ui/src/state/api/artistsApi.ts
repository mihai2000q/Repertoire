import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Artist, { BandMemberRole } from '../../types/models/Artist.ts'
import {
  AddPerfectRehearsalsToArtistsRequest,
  BulkDeleteArtistsRequest,
  DeleteArtistRequest,
  GetArtistsRequest,
  SaveImageToArtistRequest
} from '../../types/requests/ArtistRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'

const artistsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getArtists: build.query<WithTotalCountResponse<Artist>, GetArtistsRequest>({
      query: (arg) => `artists${createQueryParams(arg)}`,
      providesTags: ['Artists', 'Albums', 'Songs']
    }),
    getArtist: build.query<Artist, string>({
      query: (arg) => `artists/${arg}`,
      providesTags: ['Artists']
    }),
    addPerfectRehearsalsToArtists: build.mutation<
      HttpMessageResponse,
      AddPerfectRehearsalsToArtistsRequest
    >({
      query: (body) => ({
        url: 'artists/perfect-rehearsals',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Artists', 'Songs']
    }),
    bulkDeleteArtists: build.mutation<HttpMessageResponse, BulkDeleteArtistsRequest>({
      query: (body) => ({
        url: `artists/bulk-delete`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Artists']
    }),
    saveImageToArtist: build.mutation<HttpMessageResponse, SaveImageToArtistRequest>({
      query: (request) => ({
        url: 'artists/images',
        method: 'PUT',
        body: createFormData(request),
        formData: true
      }),
      invalidatesTags: ['Artists']
    }),
    deleteArtist: build.mutation<HttpMessageResponse, DeleteArtistRequest>({
      query: (arg) => ({
        url: `artists/${arg.id}`,
        method: 'DELETE',
        params: { ...arg, id: undefined }
      }),
      invalidatesTags: ['Artists']
    }),

    // band member - roles
    getBandMemberRoles: build.query<BandMemberRole[], void>({
      query: () => 'artists/band-members/roles',
      providesTags: ['BandMemberRoles']
    })
  })
})

export const {
  useGetArtistQuery,
  useGetArtistsQuery,
  useAddPerfectRehearsalsToArtistsMutation,
  useBulkDeleteArtistsMutation,
  useSaveImageToArtistMutation,
  useDeleteArtistMutation,
  useGetBandMemberRolesQuery
} = artistsApi
