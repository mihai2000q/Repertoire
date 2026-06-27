import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../types/models/Playlist.ts'
import {
  AddAlbumsToPlaylistRequest,
  AddArtistsToPlaylistRequest,
  AddPerfectPlaylistSongRehearsalsRequest,
  AddPerfectRehearsalsToPlaylistsRequest,
  AddSongsToPlaylistRequest,
  BulkDeletePlaylistsRequest,
  GetPlaylistRequest,
  GetPlaylistSongsRequest,
  GetPlaylistsRequest,
  SaveImageToPlaylistRequest
} from '../../types/requests/PlaylistRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'
import {
  AddAlbumsToPlaylistResponse,
  AddArtistsToPlaylistResponse,
  AddSongsToPlaylistResponse
} from '../../types/responses/PlaylistResponses.ts'
import Song from '../../types/models/Song.ts'

const playlistsApi = api.injectEndpoints({
  endpoints: (build) => ({
    // Queries
    getPlaylists: build.query<WithTotalCountResponse<Playlist>, GetPlaylistsRequest>({
      query: (arg) => `playlists${createQueryParams(arg)}`,
      providesTags: ['Playlists', 'Songs']
    }),
    getPlaylist: build.query<Playlist, GetPlaylistRequest>({
      query: (arg) => `playlists/${arg.id}${createQueryParams({ ...arg, id: undefined })}`,
      providesTags: ['Playlists']
    }),

    // Infinite queries
    getInfinitePlaylists: build.infiniteQuery<
      WithTotalCountResponse<Playlist>,
      GetPlaylistsRequest,
      { currentPage: number; pageSize: number }
    >({
      infiniteQueryOptions: {
        initialPageParam: {
          currentPage: 1,
          pageSize: 20
        },
        getNextPageParam: (lastPage, __, lastPageParam, ___, args) => {
          const pageSize = args.pageSize ?? lastPageParam.pageSize

          const totalPlaylists = lastPageParam.currentPage * pageSize
          const remainingPlaylists = lastPage?.totalCount - totalPlaylists

          if (remainingPlaylists <= 0) return undefined

          return {
            ...lastPageParam,
            currentPage: lastPageParam.currentPage + 1
          }
        }
      },
      query: ({ queryArg, pageParam }) => {
        const newQueryParams: GetPlaylistsRequest = {
          ...queryArg,
          currentPage: pageParam.currentPage,
          pageSize: queryArg.pageSize ?? pageParam.pageSize
        }
        return `playlists${createQueryParams(newQueryParams)}`
      },
      providesTags: ['Playlists', 'Songs']
    }),

    // Mutations
    addPerfectRehearsalsToPlaylists: build.mutation<
      HttpMessageResponse,
      AddPerfectRehearsalsToPlaylistsRequest
    >({
      query: (body) => ({
        url: 'playlists/perfect-rehearsals',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Playlists', 'Songs']
    }),
    bulkDeletePlaylists: build.mutation<HttpMessageResponse, BulkDeletePlaylistsRequest>({
      query: (body) => ({
        url: `playlists/bulk-delete`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Playlists']
    }),
    saveImageToPlaylist: build.mutation<HttpMessageResponse, SaveImageToPlaylistRequest>({
      query: (request) => ({
        url: 'playlists/images',
        method: 'PUT',
        body: createFormData(request),
        formData: true
      }),
      invalidatesTags: ['Playlists']
    }),
    deleteImageFromPlaylist: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `playlists/images/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Playlists']
    }),
    deletePlaylist: build.mutation<HttpMessageResponse, string>({
      query: (arg) => ({
        url: `playlists/${arg}`,
        method: 'DELETE'
      }),
      invalidatesTags: ['Playlists']
    }),

    addArtistsToPlaylist: build.mutation<AddArtistsToPlaylistResponse, AddArtistsToPlaylistRequest>(
      {
        query: (body) => ({
          url: `playlists/add-artists`,
          method: 'POST',
          body: body
        }),
        invalidatesTags: ['Songs'],
        transformResponse: (response: AddArtistsToPlaylistResponse) => ({
          ...response,
          duplicateArtistIds: response.duplicateArtistIds ?? []
        })
      }
    ),
    addAlbumsToPlaylist: build.mutation<AddAlbumsToPlaylistResponse, AddAlbumsToPlaylistRequest>({
      query: (body) => ({
        url: `playlists/add-albums`,
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs'],
      transformResponse: (response: AddAlbumsToPlaylistResponse) => ({
        ...response,
        duplicateAlbumIds: response.duplicateAlbumIds ?? []
      })
    }),

    // songs
    getInfinitePlaylistSongs: build.infiniteQuery<
      WithTotalCountResponse<Song>,
      GetPlaylistSongsRequest,
      { currentPage: number; pageSize: number }
    >({
      infiniteQueryOptions: {
        initialPageParam: {
          currentPage: 1,
          pageSize: 20
        },
        getNextPageParam: (lastPage, __, lastPageParam, ___, args) => {
          const pageSize = args.pageSize ?? lastPageParam.pageSize

          const totalSongs = lastPageParam.currentPage * pageSize
          const remainingSongs = lastPage?.totalCount - totalSongs

          if (remainingSongs <= 0) return undefined

          return {
            ...lastPageParam,
            currentPage: lastPageParam.currentPage + 1
          }
        }
      },
      query: ({ queryArg, pageParam }) => {
        const newQueryParams: GetPlaylistSongsRequest = {
          ...queryArg,
          id: undefined,
          currentPage: pageParam.currentPage,
          pageSize: queryArg.pageSize ?? pageParam.pageSize
        }
        return `playlists/songs/${queryArg.id}${createQueryParams(newQueryParams)}`
      },
      providesTags: ['Songs', 'Albums', 'Artists']
    }),
    addPerfectPlaylistSongRehearsals: build.mutation<
      HttpMessageResponse,
      AddPerfectPlaylistSongRehearsalsRequest
    >({
      query: (body) => ({
        url: 'playlists/songs/perfect-rehearsals',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    addSongsToPlaylist: build.mutation<AddSongsToPlaylistResponse, AddSongsToPlaylistRequest>({
      query: (body) => ({
        url: `playlists/songs/add`,
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    })
  })
})

export const {
  useGetPlaylistsQuery,
  useGetPlaylistQuery,
  useGetInfinitePlaylistsInfiniteQuery,
  useGetInfinitePlaylistSongsInfiniteQuery,
  useAddPerfectRehearsalsToPlaylistsMutation,
  useBulkDeletePlaylistsMutation,
  useSaveImageToPlaylistMutation,
  useDeletePlaylistMutation,
  useAddArtistsToPlaylistMutation,
  useAddAlbumsToPlaylistMutation,
  useAddSongsToPlaylistMutation,
  useAddPerfectPlaylistSongRehearsalsMutation
} = playlistsApi
