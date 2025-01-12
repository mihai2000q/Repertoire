import useShowUnknownArtist from './useShowUnknownArtist.ts'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import Album from '../types/models/Album.ts'
import Song from '../types/models/Song.ts'
import { reduxRenderHook } from '../test-utils.tsx'
import {waitFor} from "@testing-library/react";

describe('use Show Unknown Artist', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const getSongsWithoutAnArtist = http.get('/songs', async () => {
    const response: WithTotalCountResponse<Song> = { models: [], totalCount: 1 }
    return HttpResponse.json(response)
  })

  const getAlbumsWithoutAnArtist = http.get('/albums', async () => {
    const response: WithTotalCountResponse<Album> = { models: [], totalCount: 1 }
    return HttpResponse.json(response)
  })

  const getEmptySongs = http.get('/songs', async () => {
    const response: WithTotalCountResponse<Song> = { models: [], totalCount: 0 }
    return HttpResponse.json(response)
  })

  const getEmptyAlbums = http.get('/albums', async () => {
    const response: WithTotalCountResponse<Album> = { models: [], totalCount: 0 }
    return HttpResponse.json(response)
  })

  it('should return true when there are songs without an artist', async () => {
    server.use(getSongsWithoutAnArtist, getEmptyAlbums)

    const [{ result }] = reduxRenderHook(() => useShowUnknownArtist())

    await waitFor(() => expect(result.current).toBeTruthy())
  })

  it('should return true when there are albums without an artist', async () => {
    server.use(getAlbumsWithoutAnArtist, getEmptySongs)

    const [{ result }] = reduxRenderHook(() => useShowUnknownArtist())

    await waitFor(() => expect(result.current).toBeTruthy())
  })

  it('should return true when there are albums and songs without an artist', async () => {
    server.use(getAlbumsWithoutAnArtist, getSongsWithoutAnArtist)

    const [{ result }] = reduxRenderHook(() => useShowUnknownArtist())

    await waitFor(() => expect(result.current).toBeTruthy())
  })

  it('should return false when there are no songs or albums without an artist', async () => {
    server.use(getEmptySongs, getEmptyAlbums)

    const [{ result }] = reduxRenderHook(() => useShowUnknownArtist())

    await waitFor(() => expect(result.current).toBeFalsy())
  })

  it('should return false on first render', async () => {
    server.use(getSongsWithoutAnArtist, getAlbumsWithoutAnArtist)

    const [{ result }] = reduxRenderHook(() => useShowUnknownArtist())

    expect(result.current).toBeFalsy()
  })
})
