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
    // Arrange
    server.use(getSongsWithoutAnArtist, getEmptyAlbums)

    // Act
    const [{ result }] = reduxRenderHook(() => useShowUnknownArtist())

    // Assert
    await waitFor(() => expect(result.current).toBeTruthy())
  })

  it('should return true when there are albums without an artist', async () => {
    // Arrange
    server.use(getAlbumsWithoutAnArtist, getEmptySongs)

    // Act
    const [{ result }] = reduxRenderHook(() => useShowUnknownArtist())

    // Assert
    await waitFor(() => expect(result.current).toBeTruthy())
  })

  it('should return true when there are albums and songs without an artist', async () => {
    // Arrange
    server.use(getAlbumsWithoutAnArtist, getSongsWithoutAnArtist)

    // Act
    const [{ result }] = reduxRenderHook(() => useShowUnknownArtist())

    // Assert
    await waitFor(() => expect(result.current).toBeTruthy())
  })

  it('should return false when there are no songs or albums without an artist', async () => {
    // Arrange
    server.use(getEmptySongs, getEmptyAlbums)

    // Act
    const [{ result }] = reduxRenderHook(() => useShowUnknownArtist())

    // Assert
    await waitFor(() => expect(result.current).toBeFalsy())
  })

  it('should return false on first render', async () => {
    // Arrange
    server.use(getSongsWithoutAnArtist, getAlbumsWithoutAnArtist)

    // Act
    const [{ result }] = reduxRenderHook(() => useShowUnknownArtist())

    // Assert
    expect(result.current).toBeFalsy()
  })
})
