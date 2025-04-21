import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Artist, { BandMemberRole } from '../../types/models/Artist.ts'
import {
  AddAlbumsToArtistRequest,
  AddSongsToArtistRequest,
  CreateArtistRequest,
  CreateBandMemberRequest,
  DeleteArtistRequest,
  DeleteBandMemberRequest,
  GetArtistsRequest,
  MoveBandMemberRequest,
  RemoveAlbumsFromArtistRequest,
  RemoveSongsFromArtistRequest,
  SaveImageToArtistRequest,
  SaveImageToBandMemberRequest,
  UpdateArtistRequest,
  UpdateBandMemberRequest
} from '../../types/requests/ArtistRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'
import { ArtistFiltersMetadata } from '../../types/models/FiltersMetadata.ts'

const artistsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getArtists: build.query<WithTotalCountResponse<Artist>, GetArtistsRequest>({
      query: (arg) => ({
        url: `artists${createQueryParams(arg)}`
      }),
      providesTags: ['Artists', 'Albums', 'Songs']
    }),
    getArtist: build.query<Artist, string>({
      query: (arg) => `artists/${arg}`,
      providesTags: ['Artists']
    }),
    getArtistFiltersMetadata: build.query<ArtistFiltersMetadata, void>({
      query: () => 'artists/filters-metadata',
      providesTags: ['Albums', 'Artists', 'Songs']
    }),
    createArtist: build.mutation<{ id: string }, CreateArtistRequest>({
      query: (body) => ({
        url: 'artists',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Artists']
    }),
    updateArtist: build.mutation<HttpMessageResponse, UpdateArtistRequest>({
      query: (body) => ({
        url: 'artists',
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
    deleteImageFromArtist: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `artists/images/${arg}`,
        method: 'DELETE'
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

    addAlbumsToArtist: build.mutation<HttpMessageResponse, AddAlbumsToArtistRequest>({
      query: (body) => ({
        url: `artists/add-albums`,
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Artists', 'Albums']
    }),
    removeAlbumsFromArtist: build.mutation<HttpMessageResponse, RemoveAlbumsFromArtistRequest>({
      query: (body) => ({
        url: `artists/remove-albums`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Artists', 'Albums']
    }),

    addSongsToArtist: build.mutation<HttpMessageResponse, AddSongsToArtistRequest>({
      query: (body) => ({
        url: `artists/add-songs`,
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Artists', 'Songs']
    }),
    removeSongsFromArtist: build.mutation<HttpMessageResponse, RemoveSongsFromArtistRequest>({
      query: (body) => ({
        url: `artists/remove-songs`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Artists', 'Songs']
    }),

    // band member
    createBandMember: build.mutation<{ id: string }, CreateBandMemberRequest>({
      query: (body) => ({
        url: 'artists/band-members',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Artists']
    }),
    updateBandMember: build.mutation<HttpMessageResponse, UpdateBandMemberRequest>({
      query: (body) => ({
        url: 'artists/band-members',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Artists']
    }),
    moveBandMember: build.mutation<HttpMessageResponse, MoveBandMemberRequest>({
      query: (body) => ({
        url: 'artists/band-members/move',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Artists']
    }),
    saveImageToBandMember: build.mutation<HttpMessageResponse, SaveImageToBandMemberRequest>({
      query: (request) => ({
        url: 'artists/band-members/images',
        method: 'PUT',
        body: createFormData(request),
        formData: true
      }),
      invalidatesTags: ['Artists']
    }),
    deleteImageFromBandMember: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `artists/band-members/images/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Artists']
    }),
    deleteBandMember: build.mutation<HttpMessageResponse, DeleteBandMemberRequest>({
      query: (arg) => ({
        url: `artists/band-members/${arg.id}/from/${arg.artistId}`,
        method: 'DELETE'
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
  useLazyGetArtistFiltersMetadataQuery,
  useCreateArtistMutation,
  useUpdateArtistMutation,
  useSaveImageToArtistMutation,
  useDeleteImageFromArtistMutation,
  useDeleteArtistMutation,
  useAddSongsToArtistMutation,
  useRemoveSongsFromArtistMutation,
  useAddAlbumsToArtistMutation,
  useRemoveAlbumsFromArtistMutation,
  useCreateBandMemberMutation,
  useUpdateBandMemberMutation,
  useMoveBandMemberMutation,
  useSaveImageToBandMemberMutation,
  useDeleteImageFromBandMemberMutation,
  useDeleteBandMemberMutation,
  useGetBandMemberRolesQuery
} = artistsApi
