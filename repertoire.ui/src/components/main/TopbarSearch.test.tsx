import { reduxRouterRender } from '../../test-utils.tsx'
import TopbarSearch from './TopbarSearch.tsx'
import { screen, waitFor } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { http, HttpResponse, ws } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import {
  AlbumSearch,
  ArtistSearch,
  PlaylistSearch,
  SearchBase,
  SongSearch
} from '../../types/models/Search.ts'
import SearchType from '../../types/enums/SearchType.ts'
import { userEvent } from '@testing-library/user-event'
import { beforeEach } from 'vitest'

describe('Topbar Search', () => {
  const artists: ArtistSearch[] = [
    {
      id: 'artist-id',
      imageUrl: 'artist.png',
      name: 'Some Artist',
      type: SearchType.Artist
    }
  ]
  const albums: AlbumSearch[] = [
    {
      id: 'album-id',
      title: 'Some Album',
      imageUrl: 'album.png',
      type: SearchType.Album,
      artist: {
        id: 'album-artist-id',
        name: 'Some Other Artist'
      }
    }
  ]
  const songs: SongSearch[] = [
    {
      id: 'song-id',
      title: 'Some Song',
      imageUrl: 'song.png',
      type: SearchType.Song,
      artist: {
        id: 'artist-id2',
        imageUrl: 'song-artist.png',
        name: 'Some Song Artist'
      },
      album: {
        id: 'album-id2',
        imageUrl: 'song-album.png',
        title: 'Some Song Album'
      }
    }
  ]
  const playlists: PlaylistSearch[] = [
    {
      id: 'playlist-id',
      title: 'Some Playlist',
      imageUrl: 'playlist.png',
      type: SearchType.Playlist
    }
  ]

  const searchResults: SearchBase[] = [...artists, ...albums, ...songs, ...playlists]

  const handlers = [
    http.get('/centrifugo/token', async () => {
      return HttpResponse.json('')
    }),
    http.get('/search', () => {
      const result: WithTotalCountResponse<SearchBase> = {
        models: searchResults,
        totalCount: searchResults.length
      }
      return HttpResponse.json(result)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  beforeEach(() => {
    const centrifugoUrl = 'wss://chat.example.com'
    vi.stubEnv('VITE_CENTRIFUGO_URL', centrifugoUrl)
    const search = ws.link(centrifugoUrl)
    server.use(search.addEventListener('connection', () => {}))
  })

  afterEach(() => server.resetHandlers())

  afterAll(() => {
    vi.unstubAllEnvs()
    server.close()
  })

  it('should render and display recently added', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<TopbarSearch />)

    expect(screen.getByRole('searchbox', { name: /search/i })).toBeInTheDocument()

    // chips
    await user.click(screen.getByRole('searchbox', { name: /search/i }))

    expect(screen.getByText(/artists/i)).toBeInTheDocument()
    expect(screen.getByText(/albums/i)).toBeInTheDocument()
    expect(screen.getByText(/songs/i)).toBeInTheDocument()
    expect(screen.getByText(/playlists/i)).toBeInTheDocument()

    expect(screen.getByText(/recently added/i)).toBeInTheDocument()

    for (const artist of artists) {
      expect(await screen.findByRole('img', { name: artist.name })).toBeInTheDocument()
      expect(screen.getByText(artist.name)).toBeInTheDocument()
    }
    for (const album of albums) {
      expect(screen.getByRole('img', { name: album.title })).toBeInTheDocument()
      expect(screen.getByText(album.title)).toBeInTheDocument()
      expect(screen.getByText(album.artist.name)).toBeInTheDocument()
    }
    for (const song of songs) {
      expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
      expect(screen.getByText(song.title)).toBeInTheDocument()
      expect(screen.getByText(song.artist.name)).toBeInTheDocument()
    }
    for (const playlist of playlists) {
      expect(screen.getByRole('img', { name: playlist.title })).toBeInTheDocument()
      expect(screen.getByText(playlist.title)).toBeInTheDocument()
    }
  })

  it('should send search request when typing in the search box', async () => {
    const user = userEvent.setup()

    const searchValue = 'random'

    reduxRouterRender(<TopbarSearch />)

    let capturedQuery: string
    server.use(
      http.get('/search', (req) => {
        capturedQuery = new URL(req.request.url).searchParams.get('query')
        const result: WithTotalCountResponse<SearchBase> = {
          models: searchResults,
          totalCount: searchResults.length
        }
        return HttpResponse.json(result)
      })
    )

    await user.type(screen.getByRole('searchbox', { name: /search/i }), searchValue)
    await waitFor(() => expect(capturedQuery).toBe(searchValue))

    expect(screen.getByText(/search results/i)).toBeInTheDocument()
  })

  it('should send search request with type from chips', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<TopbarSearch />)

    let capturedType: SearchType
    server.use(
      http.get('/search', (req) => {
        capturedType = new URL(req.request.url).searchParams.get('type') as SearchType
        const result: WithTotalCountResponse<SearchBase> = {
          models: searchResults,
          totalCount: searchResults.length
        }
        return HttpResponse.json(result)
      })
    )

    await user.click(screen.getByRole('searchbox', { name: /search/i }))

    await user.click(screen.getByText(/artists/i))
    await waitFor(() => expect(capturedType).toBe(SearchType.Artist))

    await user.click(screen.getByText(/albums/i))
    await waitFor(() => expect(capturedType).toBe(SearchType.Album))

    await user.click(screen.getByText(/songs/i))
    await waitFor(() => expect(capturedType).toBe(SearchType.Song))

    await user.click(screen.getByText(/playlists/i))
    await waitFor(() => expect(capturedType).toBe(SearchType.Playlist))
  })

  it('should display empty messages when there are no results or nothing in the library', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<TopbarSearch />)

    server.use(
      http.get('/search', () => {
        const result: WithTotalCountResponse<SearchBase> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(result)
      })
    )

    await user.click(screen.getByRole('searchbox', { name: /search/i }))

    expect(await screen.findByText(/nothing in .* library/i)).toBeInTheDocument()
    expect(screen.queryByText(/recently added/i)).not.toBeInTheDocument()

    await user.type(screen.getByRole('searchbox', { name: /search/i }), 'r')

    await waitFor(() => expect(screen.getByText(/no results/i)).toBeInTheDocument())
    expect(screen.queryByText(/search results/i)).not.toBeInTheDocument()
  })
})
