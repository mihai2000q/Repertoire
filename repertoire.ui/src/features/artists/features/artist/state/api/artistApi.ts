import { api } from '../../../../../../state/api.ts'
import HttpMessageResponse from '../../../../../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../../../../../utils/createFormData.ts'
import {
  AddAlbumsToArtistRequest,
  AddSongsToArtistRequest,
  CreateBandMemberRequest,
  DeleteBandMemberRequest,
  MoveBandMemberRequest,
  RemoveAlbumsFromArtistRequest,
  RemoveSongsFromArtistRequest,
  SaveImageToBandMemberRequest,
  UpdateArtistRequest,
  UpdateBandMemberRequest
} from '../../types/requests/ArtistRequests.ts'

const artistsApi = api.injectEndpoints({
  endpoints: (build) => ({
    updateArtist: build.mutation<HttpMessageResponse, UpdateArtistRequest>({
      query: (body) => ({
        url: 'artists',
        method: 'PUT',
        body: body
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
    })
  })
})

export const {
  useUpdateArtistMutation,
  useDeleteImageFromArtistMutation,
  useAddSongsToArtistMutation,
  useRemoveSongsFromArtistMutation,
  useAddAlbumsToArtistMutation,
  useRemoveAlbumsFromArtistMutation,
  useCreateBandMemberMutation,
  useUpdateBandMemberMutation,
  useMoveBandMemberMutation,
  useSaveImageToBandMemberMutation,
  useDeleteImageFromBandMemberMutation,
  useDeleteBandMemberMutation
} = artistsApi
