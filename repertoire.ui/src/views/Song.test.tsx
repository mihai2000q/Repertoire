import { reduxMemoryRouterRender } from '../test-utils.tsx'
import Song from './Song.tsx'
import { default as SongType } from './../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { screen } from '@testing-library/react'
import { RootState } from '../state/store.ts'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import Artist from '../types/models/Artist.ts'
import Album from '../types/models/Album.ts'

describe('Song', () => {
  const song: SongType = {
    id: '1',
    title: 'Song 1',
    description: '',
    isRecorded: false,
    rehearsals: 0,
    confidence: 0,
    progress: 0,
    sections: [],
    createdAt: '',
    updatedAt: '',
    releaseDate: null
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
    http.get('/artists', async () => {
      const response: WithTotalCountResponse<Artist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = { models: [], totalCount: 0 }
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
