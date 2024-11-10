import { api } from './api'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse'
import Artist from "../types/models/Artist.ts";
import {
  CreateArtistRequest,
  GetArtistsRequest,
  SaveImageToArtistRequest,
  UpdateArtistRequest
} from "../types/requests/ArtistRequests.ts";
import HttpMessageResponse from "../types/responses/HttpMessageResponse.ts";
import createFormData from "../utils/createFormData.ts";

const songsApi = api.injectEndpoints({
  endpoints: (build) => ({
    getArtists: build.query<WithTotalCountResponse<Artist>, GetArtistsRequest>({
      query: (arg) => ({
        url: 'artists',
        params: arg
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
    })
  })
})

export const { useGetArtistsQuery } = songsApi
