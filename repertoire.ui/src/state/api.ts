import { createApi } from '@reduxjs/toolkit/query/react'
import { queryWithRedirection } from './api.query'

export const api = createApi({
  baseQuery: queryWithRedirection,
  reducerPath: 'api',
  tagTypes: [
    'Songs',
    'Albums',
    'Artists',
    'Playlists',
    'GuitarTunings',
    'Instruments',
    'BandMemberRoles',
    'SongSectionTypes',
    'User',
    'Search'
  ],
  endpoints: () => ({})
})
