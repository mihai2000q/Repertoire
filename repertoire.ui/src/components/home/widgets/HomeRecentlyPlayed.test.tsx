import { emptyAlbum, emptyArtist, emptySong, reduxRouterRender } from '../../../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import Song from '../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import HomeRecentlyPlayed from './HomeRecentlyPlayed.tsx'
import { expect } from 'vitest'
import { userEvent } from '@testing-library/user-event'
import Artist from '../../../types/models/Artist.ts'
import Album from '../../../types/models/Album.ts'

describe('Home Recently Played', () => {
  const artists: Artist[] = [
    {
      ...emptyArtist,
      id: '1',
      name: 'Artist 1',
      progress: 10,
      lastTimePlayed: '2025-01-22T10:30'
    },
    {
      ...emptyArtist,
      id: '2',
      name: 'Artist 2',
      progress: 10,
      lastTimePlayed: '2025-01-22T12:22'
    }
  ]

  const albums: Album[] = [
    {
      ...emptyAlbum,
      id: '1',
      title: 'Album 1',
      progress: 10,
      lastTimePlayed: '2025-01-22T10:30'
    },
    {
      ...emptyAlbum,
      id: '2',
      title: 'Album 2',
      progress: 10,
      lastTimePlayed: '2025-01-22T12:22'
    }
  ]

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
    http.get('/artists', async () => {
      const response: WithTotalCountResponse<Artist> = {
        models: artists,
        totalCount: artists.length
      }
      return HttpResponse.json(response)
    }),
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = {
        models: albums,
        totalCount: albums.length
      }
      return HttpResponse.json(response)
    }),
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

  afterEach(() => {
    server.resetHandlers()
    localStorage.clear()
  })

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<HomeRecentlyPlayed />)

    expect(screen.getByText(/recently played/i)).toBeInTheDocument()
    // early check due to the below await hover
    expect(screen.getByTestId('recently-played-loader')).toBeInTheDocument()

    await user.hover(screen.getByLabelText('recently-played-widget'))
    expect(
      await screen.findByRole('button', { name: 'recently-played-artists' })
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'recently-played-albums' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'recently-played-songs' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'recently-played-songs' })).toHaveAttribute(
      'aria-selected',
      'true'
    )

    expect(screen.getByText(/title/i)).toBeInTheDocument()
    expect(screen.getByText(/progress/i)).toBeInTheDocument()
    expect(screen.getByLabelText('last-time-played-icon')).toBeInTheDocument()

    // default - songs
    for (const song of songs) {
      expect(await screen.findByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    }

    // albums
    await user.click(screen.getByRole('button', { name: 'recently-played-albums' }))
    expect(screen.getByRole('button', { name: 'recently-played-albums' })).toHaveAttribute(
      'aria-selected',
      'true'
    )
    for (const album of albums) {
      expect(screen.getByLabelText(`album-card-${album.title}`)).toBeInTheDocument()
    }

    // artists
    await user.click(screen.getByRole('button', { name: 'recently-played-artists' }))
    expect(screen.getByRole('button', { name: 'recently-played-artists' })).toHaveAttribute(
      'aria-selected',
      'true'
    )

    expect(screen.queryByText(/title/i)).not.toBeInTheDocument()
    expect(screen.getByText(/name/i)).toBeInTheDocument()
    for (const artist of artists) {
      expect(screen.getByLabelText(`artist-card-${artist.name}`)).toBeInTheDocument()
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

  it('should display empty message when there are no albums', async () => {
    const user = userEvent.setup()

    server.use(
      http.get('/albums', async () => {
        const response: WithTotalCountResponse<Album> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<HomeRecentlyPlayed />)

    await user.click(screen.getByRole('button', { name: 'recently-played-albums' }))

    expect(screen.getByText(/recently played/i)).toBeInTheDocument()

    expect(await screen.findByText(/no albums/i)).toBeInTheDocument()

    expect(screen.queryByText(/title/i)).not.toBeInTheDocument()
    expect(screen.queryByText(/progress/i)).not.toBeInTheDocument()
    expect(screen.queryByLabelText('last-time-played-icon')).not.toBeInTheDocument()
  })

  it('should display empty message when there are no artists', async () => {
    const user = userEvent.setup()

    server.use(
      http.get('/artists', async () => {
        const response: WithTotalCountResponse<Artist> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<HomeRecentlyPlayed />)

    await user.click(screen.getByRole('button', { name: 'recently-played-artists' }))

    expect(screen.getByText(/recently played/i)).toBeInTheDocument()

    expect(await screen.findByText(/no artists/i)).toBeInTheDocument()

    expect(screen.queryByText(/title/i)).not.toBeInTheDocument()
    expect(screen.queryByText(/progress/i)).not.toBeInTheDocument()
    expect(screen.queryByLabelText('last-time-played-icon')).not.toBeInTheDocument()
  })
})
