import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import Song from '../types/models/Song.ts'
import { reduxRenderHook } from '../test-utils.tsx'
import { waitFor } from '@testing-library/react'
import useShowUnknownAlbum from './useShowUnknownAlbum.ts'

describe('use Show Unknown Album', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const getSongsWithoutAnAlbum = http.get('/songs', async () => {
    const response: WithTotalCountResponse<Song> = { models: [], totalCount: 1 }
    return HttpResponse.json(response)
  })

  const getEmptySongs = http.get('/songs', async () => {
    const response: WithTotalCountResponse<Song> = { models: [], totalCount: 0 }
    return HttpResponse.json(response)
  })

  it('should return true when there are songs without an album', async () => {
    server.use(getSongsWithoutAnAlbum)

    const [{ result }] = reduxRenderHook(() => useShowUnknownAlbum())

    await waitFor(() => expect(result.current).toBeTruthy())
  })

  it('should return false when there are no songs without an album', async () => {
    server.use(getEmptySongs)

    const [{ result }] = reduxRenderHook(() => useShowUnknownAlbum())

    await waitFor(() => expect(result.current).toBeFalsy())
  })

  it('should return false on first render', async () => {
    server.use(getSongsWithoutAnAlbum)

    const [{ result }] = reduxRenderHook(() => useShowUnknownAlbum())

    expect(result.current).toBeFalsy()
  })
})
