import { reduxMemoryRouterRender } from '../test-utils.tsx'
import Song from './Song.tsx'
import { default as SongType } from './../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { screen } from '@testing-library/react'

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
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxMemoryRouterRender(<Song />, '/song/:id', [`/song/${song.id}`])

    expect(screen.getByTestId('song-loader')).toBeInTheDocument()
    expect(await screen.findByLabelText('header-panel-card')).toBeInTheDocument()
    expect(screen.getByLabelText('song-information-card')).toBeInTheDocument()
    expect(screen.getByLabelText('song-overall-card')).toBeInTheDocument()
    expect(screen.getByLabelText('song-links-card')).toBeInTheDocument()
    expect(screen.getByLabelText('song-description-card')).toBeInTheDocument()
    expect(screen.getByLabelText('song-sections')).toBeInTheDocument()
  })
})
