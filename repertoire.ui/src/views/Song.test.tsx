import { emptySong, reduxMemoryRouterRender } from '../test-utils.tsx'
import Song from './Song.tsx'
import { default as SongType } from './../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { screen } from '@testing-library/react'
import { RootState } from '../state/store.ts'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { SearchBase } from '../types/models/Search.ts'

describe('Song', () => {
  const song: SongType = {
    ...emptySong,
    id: '1',
    title: 'Song 1'
  }

  const handlers = [
    http.get(`/songs/${song.id}`, async () => {
      return HttpResponse.json(song)
    }),
    http.get(`/songs/sections/types`, () => {
      return HttpResponse.json([])
    }),
    http.get(`/songs/guitar-tunings`, () => {
      return HttpResponse.json([])
    }),
    http.get(`/songs/instruments`, () => {
      return HttpResponse.json([])
    }),
    http.get('/search', async () => {
      const response: WithTotalCountResponse<SearchBase> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const [_, store] = reduxMemoryRouterRender(<Song />, '/song/:id', [`/song/${song.id}`])

    expect(screen.getByTestId('song-loader')).toBeInTheDocument()
    expect(await screen.findByLabelText('header-panel-card')).toBeInTheDocument()
    expect(screen.getByLabelText('song-information-card')).toBeInTheDocument()
    expect(screen.getByLabelText('song-overall-card')).toBeInTheDocument()
    expect(screen.getByLabelText('song-links-card')).toBeInTheDocument()
    expect(screen.getByLabelText('song-description-card')).toBeInTheDocument()
    expect(screen.getByLabelText('song-sections')).toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toBe(song.title)
  })
})
