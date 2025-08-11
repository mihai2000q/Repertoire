import { emptySong, reduxRouterRender } from '../../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Song from '../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import HomeRecentlyPlayed from './HomeRecentlyPlayed.tsx'
import { expect } from 'vitest'

describe('Home Recently Played', () => {
  const songs: Song[] = [
    {
      ...emptySong,
      id: '1',
      title: 'Song 1',
      progress: 10,
      lastTimePlayed: '2025-01-22T10:30'
    },
    {
      ...emptySong,
      id: '2',
      title: 'Song 2',
      progress: 10,
      lastTimePlayed: '2025-01-22T12:22'
    }
  ]

  const handlers = [
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
        models: songs,
        totalCount: songs.length
      }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRouterRender(<HomeRecentlyPlayed />)

    expect(screen.getByText(/recently played/i)).toBeInTheDocument()

    expect(screen.getByText(/title/i)).toBeInTheDocument()
    expect(screen.getByText(/progress/i)).toBeInTheDocument()
    expect(screen.getByLabelText('last-time-played-icon')).toBeInTheDocument()

    for (const song of songs) {
      expect(await screen.findByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    }
  })

  it('should display empty message when there are no songs', async () => {
    server.use(
      http.get('/songs', async () => {
        const response: WithTotalCountResponse<Song> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<HomeRecentlyPlayed />)

    expect(screen.getByText(/recently played/i)).toBeInTheDocument()

    expect(await screen.findByText(/no songs/i)).toBeInTheDocument()

    expect(screen.queryByText(/title/i)).not.toBeInTheDocument()
    expect(screen.queryByText(/progress/i)).not.toBeInTheDocument()
    expect(screen.queryByLabelText('last-time-played-icon')).not.toBeInTheDocument()
  })
})
