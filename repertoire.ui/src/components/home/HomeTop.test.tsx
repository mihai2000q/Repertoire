import { emptyAlbum, emptyArtist, emptySong, reduxRender } from '../../test-utils.tsx'
import HomeTop from './HomeTop.tsx'
import { screen } from '@testing-library/react'
import Album from '../../types/models/Album.ts'
import Artist from '../../types/models/Artist.ts'
import Song from '../../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'

describe('Home Top', () => {
  const albums: Album[] = [
    {
      ...emptyAlbum,
      id: '1',
      title: 'Album 1',
      imageUrl: 'something-album.png',
      artist: emptyArtist
    },
    {
      ...emptyAlbum,
      id: '2',
      title: 'Album 2'
    }
  ]

  const artists: Artist[] = [
    {
      ...emptyArtist,
      id: '1',
      name: 'Artist 1'
    },
    {
      ...emptyArtist,
      id: '2',
      name: 'Artist 2'
    }
  ]

  const songs: Song[] = [
    {
      ...emptySong,
      id: '1',
      title: 'Song 1',
      imageUrl: 'something.png',
      album: {
        ...emptyAlbum,
        imageUrl: 'something-album.png'
      }
    },
    {
      ...emptySong,
      id: '2',
      title: 'Song 2',
      album: {
        ...emptyAlbum,
        imageUrl: 'something-album.png'
      },
      artist: {
        ...emptyArtist,
        id: '1',
        name: 'Artist'
      }
    },
    {
      ...emptySong,
      id: '3',
      title: 'Song 11',
      imageUrl: 'something.png',
      album: emptyAlbum
    },
    {
      ...emptySong,
      id: '4',
      title: 'Song 12',
      album: emptyAlbum
    },
    {
      ...emptySong,
      id: '5',
      title: 'Song 512',
      imageUrl: 'something.png'
    },
    {
      ...emptySong,
      id: '6',
      title: 'Song 6'
    }
  ]

  const handlers = [
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
        models: songs,
        totalCount: songs.length
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
    http.get('/artists', async () => {
      const response: WithTotalCountResponse<Artist> = {
        models: artists,
        totalCount: artists.length
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
    reduxRender(<HomeTop />)

    expect(screen.getByText(/welcome back/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'back' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'forward' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /songs/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /albums/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /albums/i })).toHaveAttribute('aria-selected', 'true')
    expect(screen.getByRole('button', { name: /artists/i })).toBeInTheDocument()

    for (const album of albums) {
      expect(await screen.findByLabelText(`album-card-${album.title}`)).toBeInTheDocument()
    }
  })

  it('should display songs, artists and albums when switching tabs', async () => {
    const user = userEvent.setup()

    reduxRender(<HomeTop />)

    const songsTab = screen.getByRole('button', { name: /songs/i })
    const albumsTab = screen.getByRole('button', { name: /albums/i })
    const artistsTab = screen.getByRole('button', { name: /artists/i })

    // songs' tab
    await user.click(songsTab)
    expect(songsTab).toHaveAttribute('aria-selected', 'true')

    for (const song of songs) {
      expect(await screen.findByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    }

    // albums' tab
    await user.click(albumsTab)
    expect(albumsTab).toHaveAttribute('aria-selected', 'true')

    for (const album of albums) {
      expect(await screen.findByLabelText(`album-card-${album.title}`)).toBeInTheDocument()
    }

    // artists' tab
    await user.click(artistsTab)
    expect(artistsTab).toHaveAttribute('aria-selected', 'true')

    for (const artist of artists) {
      expect(await screen.findByLabelText(`artist-card-${artist.name}`)).toBeInTheDocument()
    }
  })

  it('should display empty message when there are no songs, albums or artists', async () => {
    const user = userEvent.setup()

    server.use(
      http.get('/songs', async () => {
        const response: WithTotalCountResponse<Song> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      }),
      http.get('/albums', async () => {
        const response: WithTotalCountResponse<Album> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      }),
      http.get('/artists', async () => {
        const response: WithTotalCountResponse<Artist> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRender(<HomeTop />)

    expect(await screen.findByText(/no albums/i)).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: /songs/i }))
    expect(await screen.findByText(/no songs/i)).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: /artists/i }))
    expect(await screen.findByText(/no artists/i)).toBeInTheDocument()
  })
})
