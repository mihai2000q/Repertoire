import { emptyAlbum, emptyArtist, emptySong, reduxRouterRender } from '../../test-utils.tsx'
import HomeTop from './HomeTop.tsx'
import { screen, waitFor } from '@testing-library/react'
import Album from '../../types/models/Album.ts'
import Artist from '../../types/models/Artist.ts'
import Song from '../../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import homeTopOrderEntities from '../../data/home/homeTopOrderEntities.ts'
import HomeTopEntity from '../../types/enums/HomeTopEntity.ts'
import OrderType from '../../types/enums/OrderType.ts'

describe('Home Top', () => {
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
    reduxRouterRender(<HomeTop />)

    expect(screen.getByText(/welcome back/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /artists/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /albums/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /albums/i })).toHaveAttribute('aria-selected', 'true')
    expect(screen.getByRole('button', { name: /songs/i })).toBeInTheDocument()

    expect(screen.getByRole('button', { name: 'order' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'order' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'back' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'forward' })).toBeInTheDocument()

    for (const album of albums) {
      expect(await screen.findByLabelText(`album-card-${album.title}`)).toBeInTheDocument()
    }
    expect(screen.getByRole('button', { name: 'order' })).not.toBeDisabled()
  })

  it('should display artists, albums and songs when switching tabs', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<HomeTop />)

    const artistsTab = screen.getByRole('button', { name: /artists/i })
    const albumsTab = screen.getByRole('button', { name: /albums/i })
    const songsTab = screen.getByRole('button', { name: /songs/i })

    // artists' tab
    await user.click(artistsTab)
    expect(artistsTab).toHaveAttribute('aria-selected', 'true')

    for (const artist of artists) {
      expect(await screen.findByLabelText(`artist-card-${artist.name}`)).toBeInTheDocument()
    }

    // albums' tab
    await user.click(albumsTab)
    expect(albumsTab).toHaveAttribute('aria-selected', 'true')

    for (const album of albums) {
      expect(await screen.findByLabelText(`album-card-${album.title}`)).toBeInTheDocument()
    }

    // songs' tab
    await user.click(songsTab)
    expect(songsTab).toHaveAttribute('aria-selected', 'true')

    for (const song of songs) {
      expect(await screen.findByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    }
  })

  it('should display empty message when there are no artists, albums or songs and disable the order button', async () => {
    const user = userEvent.setup()

    server.use(
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
      }),
      http.get('/songs', async () => {
        const response: WithTotalCountResponse<Song> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<HomeTop />)

    expect(await screen.findByText(/no albums/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'order' })).toBeDisabled()

    await user.click(screen.getByRole('button', { name: /songs/i }))
    expect(await screen.findByText(/no songs/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'order' })).toBeDisabled()

    await user.click(screen.getByRole('button', { name: /artists/i }))
    expect(await screen.findByText(/no artists/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'order' })).toBeDisabled()
  })

  describe('should have a default order and have the option to change order', () => {
    it('should display artists with default order and be able to change it', async () => {
      await test(HomeTopEntity.Artists, 'artists', albums)
    })

    it('should display albums with default order and be able to change it', async () => {
      await test(HomeTopEntity.Albums, 'albums', albums, true)
    })

    it('should display songs with default order and be able to change it', async () => {
      await test(HomeTopEntity.Songs, 'songs', songs)
    })

    async function test<TEntity>(
      topEntity: HomeTopEntity,
      path: string,
      entities: TEntity[],
      isDefaultTab: boolean = false
    ) {
      const user = userEvent.setup()

      const initialOrder = homeTopOrderEntities.get(topEntity)[0]
      const newOrder = homeTopOrderEntities.get(topEntity)[1]

      let capturedOrderBy: string[]
      server.use(
        http.get('/' + path, async (req) => {
          capturedOrderBy = new URL(req.request.url).searchParams.getAll('orderBy')
          const response: WithTotalCountResponse<TEntity> = {
            models: entities,
            totalCount: entities.length
          }
          return HttpResponse.json(response)
        })
      )

      reduxRouterRender(<HomeTop />)

      if (!isDefaultTab)
        await user.click(screen.getByRole('button', { name: new RegExp(path, 'i') }))

      await waitFor(() =>
        expect(capturedOrderBy).toStrictEqual([
          initialOrder.property + ' ' + initialOrder.type,
          initialOrder.thenBy[0].property + ' ' + OrderType.Ascending
        ])
      )

      await user.click(screen.getByRole('button', { name: 'order' }))
      await user.click(await screen.findByRole('menuitem', { name: newOrder.label }))

      await waitFor(() =>
        expect(capturedOrderBy).toStrictEqual([
          newOrder.property + ' ' + newOrder.type,
          newOrder.thenBy[0].property + ' ' + OrderType.Ascending
        ])
      )
    }
  })
})
