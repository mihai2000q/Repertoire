import { api } from './api'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse'
import Artist from '../types/models/Artist.ts'
import {
  AddAlbumsToArtistRequest,
  AddSongsToArtistRequest,
  CreateArtistRequest,
  GetArtistsRequest,
  RemoveAlbumsFromAristRequest,
  RemoveSongsFromArtistRequest,
  SaveImageToArtistRequest,
  UpdateArtistRequest
} from '../types/requests/ArtistRequests.ts'
import HttpMessageResponse from '../types/responses/HttpMessageResponse.ts'
import createFormData from '../utils/createFormData.ts'
import createQueryParams from "../utils/createQueryParams.ts";

const artistsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getArtists: build.query<WithTotalCountResponse<Artist>, GetArtistsRequest>({
      query: (arg) => ({
        url: `artists${createQueryParams(arg)}`,
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
    deleteArtist: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `artists/${arg}`,
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
    removeAlbumsFromArtist: build.mutation<HttpMessageResponse, RemoveAlbumsFromAristRequest>({
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
    })
  })
})

export const {
  useGetArtistQuery,
  useGetArtistsQuery,
  useCreateArtistMutation,
  useUpdateArtistMutation,
  useSaveImageToArtistMutation,
  useDeleteArtistMutation,
  useAddSongsToArtistMutation,
  useRemoveSongsFromArtistMutation,
  useAddAlbumsToArtistMutation,
  useRemoveAlbumsFromArtistMutation
} = artistsApi
