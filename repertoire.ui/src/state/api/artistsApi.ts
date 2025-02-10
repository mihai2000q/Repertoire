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

const artistsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getArtists: build.query<WithTotalCountResponse<Artist>, GetArtistsRequest>({
      query: (arg) => ({
        url: `artists${createQueryParams(arg)}`
      }),
      providesTags: ['Artists']
    }),
    getArtist: build.query<Artist, string>({
      query: (arg) => `artists/${arg}`,
      providesTags: ['Artists']
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
      invalidatesTags: ['Artists', 'Songs', 'Albums']
    }),
    saveImageToArtist: build.mutation<HttpMessageResponse, SaveImageToArtistRequest>({
      query: (request) => ({
        url: 'artists/images',
        method: 'PUT',
        body: createFormData(request),
        formData: true
      }),
      invalidatesTags: ['Artists', 'Songs', 'Albums']
    }),
    deleteImageFromArtist: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `artists/images/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Artists', 'Songs', 'Albums']
    }),
    deleteArtist: build.mutation<HttpMessageResponse, DeleteArtistRequest>({
      query: (arg) => ({
        url: `artists/${arg.id}`,
        method: 'DELETE',
        params: { ...arg, id: undefined }
      }),
      invalidatesTags: ['Artists', 'Songs', 'Albums']
    }),

    addAlbumsToArtist: build.mutation<HttpMessageResponse, AddAlbumsToArtistRequest>({
      query: (body) => ({
        url: `artists/add-albums`,
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Artists', 'Albums', 'Songs']
    }),
    removeAlbumsFromArtist: build.mutation<HttpMessageResponse, RemoveAlbumsFromArtistRequest>({
      query: (body) => ({
        url: `artists/remove-albums`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Artists', 'Albums', 'Songs']
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
      invalidatesTags: ['Artists', 'Songs']
    }),
    updateBandMember: build.mutation<HttpMessageResponse, UpdateBandMemberRequest>({
      query: (body) => ({
        url: 'artists/band-members',
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Artists', 'Songs']
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
      invalidatesTags: ['Artists', 'Songs']
    }),
    deleteImageFromBandMember: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `artists/band-members/images/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Artists', 'Songs']
    }),
    deleteBandMember: build.mutation<HttpMessageResponse, DeleteBandMemberRequest>({
      query: (arg) => ({
        url: `artists/band-members/${arg.id}/from/${arg.artistId}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Artists', 'Songs']
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
