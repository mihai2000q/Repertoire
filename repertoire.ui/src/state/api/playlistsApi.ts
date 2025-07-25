import { api } from '../api.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../types/models/Playlist.ts'
import {
  AddAlbumsToPlaylistRequest,
  AddArtistsToPlaylistRequest,
  AddSongsToPlaylistRequest,
  CreatePlaylistRequest,
  GetPlaylistRequest,
  GetPlaylistSongsRequest,
  GetPlaylistsRequest,
  MoveSongFromPlaylistRequest,
  RemoveSongsFromPlaylistRequest,
  SaveImageToPlaylistRequest,
  ShufflePlaylistSongsRequest,
  UpdatePlaylistRequest
} from '../../types/requests/PlaylistRequests.ts'
import HttpMessageResponse from '../../types/responses/HttpMessageResponse.ts'
import createFormData from '../../utils/createFormData.ts'
import createQueryParams from '../../utils/createQueryParams.ts'
import { PlaylistFiltersMetadata } from '../../types/models/FiltersMetadata.ts'
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
    getPlaylistFiltersMetadata: build.query<PlaylistFiltersMetadata, { searchBy?: string[] }>({
      query: (arg) => `playlists/filters-metadata${createQueryParams(arg)}`,
      providesTags: ['Playlists', 'Songs']
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
    createPlaylist: build.mutation<{ id: string }, CreatePlaylistRequest>({
      query: (body) => ({
        url: 'playlists',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Playlists']
    }),
    updatePlaylist: build.mutation<HttpMessageResponse, UpdatePlaylistRequest>({
      query: (body) => ({
        url: 'playlists',
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
    addSongsToPlaylist: build.mutation<AddSongsToPlaylistResponse, AddSongsToPlaylistRequest>({
      query: (body) => ({
        url: `playlists/songs/add`,
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    shufflePlaylist: build.mutation<HttpMessageResponse, ShufflePlaylistSongsRequest>({
      query: (body) => ({
        url: 'playlists/songs/shuffle',
        method: 'POST',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    moveSongFromPlaylist: build.mutation<HttpMessageResponse, MoveSongFromPlaylistRequest>({
      query: (body) => ({
        url: `playlists/songs/move`,
        method: 'PUT',
        body: body
      }),
      invalidatesTags: ['Songs']
    }),
    removeSongsFromPlaylist: build.mutation<HttpMessageResponse, RemoveSongsFromPlaylistRequest>({
      query: (body) => ({
        url: `playlists/songs/remove`,
        method: 'PUT',
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
  useGetPlaylistFiltersMetadataQuery,
  useLazyGetPlaylistFiltersMetadataQuery,
  useCreatePlaylistMutation,
  useUpdatePlaylistMutation,
  useSaveImageToPlaylistMutation,
  useDeleteImageFromPlaylistMutation,
  useDeletePlaylistMutation,
  useAddArtistsToPlaylistMutation,
  useAddAlbumsToPlaylistMutation,
  useAddSongsToPlaylistMutation,
  useMoveSongFromPlaylistMutation,
  useRemoveSongsFromPlaylistMutation
} = playlistsApi
