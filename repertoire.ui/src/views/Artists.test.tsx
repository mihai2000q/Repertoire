import { screen, waitFor } from '@testing-library/react'
import Artists from './Artists.tsx'
import {
  defaultArtistFiltersMetadata,
  emptyAlbum,
  emptyArtist,
  emptySong,
  reduxRouterRender
} from '../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { userEvent } from '@testing-library/user-event'
import Artist from '../types/models/Artist.ts'
import Song from '../types/models/Song.ts'
import Album from '../types/models/Album.ts'
import { RootState } from '../state/store.ts'
import artistsOrders from '../data/artists/artistsOrders.ts'
import FilterOperator from '../types/enums/FilterOperator.ts'
import ArtistProperty from '../types/enums/properties/ArtistProperty.ts'
import OrderType from '../types/enums/OrderType.ts'
import Playlist from '../types/models/Playlist.ts'

describe('Artists', () => {
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
  const totalCount = 2

  const initialCurrentPage = 1
  const pageSize = 40

  const getSongsWithoutArtists = () =>
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
        models: [
          {
            ...emptySong,
            title: 'Some song'
          }
        ],
        totalCount: 1
      }
      return HttpResponse.json(response)
    })

  const getAlbumsWithoutArtists = () =>
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = {
        models: [
          {
            ...emptyAlbum,
            title: 'Some Album'
          }
        ],
        totalCount: 1
      }
      return HttpResponse.json(response)
    })

  const getArtistsWithPagination = (totalCount: number = 50) =>
    http.get('/artists', async (req) => {
      const currentPage = new URL(req.request.url).searchParams.get('currentPage')
      const response: WithTotalCountResponse<Artist> =
        parseInt(currentPage) === 1
          ? {
              models: Array.from({ length: pageSize }).map((_, i) => ({
                ...artists[0],
                id: i.toString()
              })),
              totalCount: totalCount
            }
          : {
              models: Array.from({ length: totalCount - pageSize }).map((_, i) => ({
                ...artists[0],
                id: i.toString()
              })),
              totalCount: totalCount
            }
      return HttpResponse.json(response)
    })

  const getNoArtists = () =>
    http.get('/artists', async () => {
      const response: WithTotalCountResponse<Artist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    })

  const handlers = [
    http.get('/artists', async () => {
      const response: WithTotalCountResponse<Artist> = {
        models: artists,
        totalCount: totalCount
      }
      return HttpResponse.json(response)
    }),
    http.get('/playlists', async () => {
      const response: WithTotalCountResponse<Playlist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/artists/filters-metadata', async () => {
      return HttpResponse.json(defaultArtistFiltersMetadata)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => {
    server.listen()
    vi.mock('../context/MainScrollContext.tsx', () => ({
      useMainScroll: vi.fn(() => ({
        ref: { current: document.createElement('div') }
      }))
    }))
  })

  afterEach(() => {
    window.location.search = ''
    server.resetHandlers()
  })

  afterAll(() => {
    server.close()
    vi.clearAllMocks()
  })

  it('should render and display relevant info when there are artists', async () => {
    const [_, store] = reduxRouterRender(<Artists />)

    expect((store.getState() as RootState).global.documentTitle).toMatch(/artists/i)
    expect(screen.getByRole('heading', { name: /artists/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /new-artist/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-artists/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-artists/i })).toBeDisabled()
    expect(screen.getByRole('button', { name: /filter-artists/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /filter-artists/i })).toBeDisabled()

    expect(await screen.findByLabelText('new-artist-card')).toBeInTheDocument()
    expect(screen.getAllByLabelText(/artist-card-/)).toHaveLength(artists.length)
    expect(
      screen.getByText(`${initialCurrentPage} - ${artists.length} artists out of ${totalCount}`)
    ).toBeInTheDocument()
    artists.forEach((artist) =>
      expect(screen.getByLabelText(`artist-card-${artist.name}`)).toBeInTheDocument()
    )
    expect(screen.getByTestId('artists-pagination')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-artists/i })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: /filter-artists/i })).not.toBeDisabled()
  })

  it('should open the add new artist modal when clicking the new artist button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<Artists />)

    const newAlbumButton = screen.getByRole('button', { name: /new-artist/i })
    await user.click(newAlbumButton)

    expect(await screen.findByRole('heading', { name: /add new artist/i })).toBeInTheDocument()
  })

  it('should open the add new artist modal when clicking the new artist card button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<Artists />)

    const newAlbumCardButton = await screen.findByLabelText('new-artist-card')
    await user.click(newAlbumCardButton)

    expect(await screen.findByRole('heading', { name: /add new artist/i })).toBeInTheDocument()
  })

  it('should display not display some info when there are no artists', async () => {
    reduxRouterRender(<Artists />)

    server.use(getNoArtists())

    expect(await screen.findByText(/no artists/)).toBeInTheDocument()
    expect(screen.queryByLabelText('new-artist-card')).not.toBeInTheDocument()
    expect(screen.queryByLabelText('unknown-artist-card')).not.toBeInTheDocument()
    expect(screen.queryAllByLabelText(/artist-card-/)).toHaveLength(0)
    expect(screen.queryByTestId('artists-pagination')).not.toBeInTheDocument()
  })

  it('should display pagination and new artist card when the unknown artist card is shown and there are no artists', async () => {
    reduxRouterRender(<Artists />)

    server.use(getNoArtists(), getSongsWithoutArtists())

    expect(await screen.findByLabelText('unknown-artist-card')).toBeInTheDocument()
    expect(screen.queryByText(/no artists/)).not.toBeInTheDocument()
    expect(screen.queryByLabelText('new-artist-card')).toBeInTheDocument()
    expect(screen.queryAllByLabelText(/artist-card-/)).toHaveLength(0)
    expect(screen.queryByTestId('artists-pagination')).toBeInTheDocument()
  })

  it('should paginate the artists', async () => {
    const user = userEvent.setup()

    const totalCount = 50

    server.use(getArtistsWithPagination(totalCount))

    reduxRouterRender(<Artists />)

    expect(await screen.findByTestId('artists-pagination')).toBeInTheDocument()
    expect(screen.queryAllByLabelText(/artist-card-/)).toHaveLength(pageSize)
    expect(
      screen.getByText(`${initialCurrentPage} - ${pageSize} artists out of ${totalCount}`)
    ).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '1' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '2' })).toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('artists-pagination')).toBeInTheDocument()
    expect(screen.queryAllByLabelText(/artist-card-/)).toHaveLength(totalCount - pageSize)
    expect(
      screen.getByText(`${pageSize + 1} - ${totalCount} artists out of ${totalCount}`)
    ).toBeInTheDocument()
    expect(window.location.search).toContain('p=2')
  })

  it('should display the new artist card on last page, but not on the first page', async () => {
    const user = userEvent.setup()

    server.use(getArtistsWithPagination())

    reduxRouterRender(<Artists />)

    expect(await screen.findByTestId('artists-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('new-artist-card')).not.toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('artists-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('new-artist-card')).toBeInTheDocument()
  })

  it.each([
    [false, true],
    [true, false],
    [true, true]
  ])(
    'should display unknown artist card when there are songs or albums without artist',
    async (withAlbum, withSong) => {
      if (withAlbum) {
        server.use(getAlbumsWithoutArtists())
      }
      if (withSong) {
        server.use(getSongsWithoutArtists())
      }

      reduxRouterRender(<Artists />)

      expect(await screen.findByTestId('artists-pagination')).toBeInTheDocument()
      expect(screen.queryByLabelText('unknown-artist-card')).toBeInTheDocument()
      expect(
        screen.getByText(
          `${initialCurrentPage} - ${artists.length + 1} artists out of ${totalCount + 1}`
        )
      ).toBeInTheDocument()
    }
  )

  it.each([
    [false, true],
    [true, false],
    [true, true]
  ])(
    'should display the new artist and unknown artist cards on the last page, but not the first page',
    async (withAlbum, withSong) => {
      const user = userEvent.setup()

      const totalCount = 50

      server.use(getArtistsWithPagination(totalCount))
      if (withAlbum) {
        server.use(getAlbumsWithoutArtists())
      }
      if (withSong) {
        server.use(getSongsWithoutArtists())
      }

      reduxRouterRender(<Artists />)

      expect(await screen.findByTestId('artists-pagination')).toBeInTheDocument()
      expect(screen.queryByLabelText('unknown-artist-card')).not.toBeInTheDocument()
      expect(screen.queryByLabelText('new-artist-card')).not.toBeInTheDocument()
      expect(
        screen.getByText(`${initialCurrentPage} - ${pageSize} artists out of ${totalCount + 1}`)
      ).toBeInTheDocument()

      const pagination = screen.getByRole('button', { name: '2' })
      await user.click(pagination)

      expect(await screen.findByTestId('artists-pagination')).toBeInTheDocument()
      expect(screen.getByLabelText('unknown-artist-card')).toBeInTheDocument()
      expect(screen.getByLabelText('new-artist-card')).toBeInTheDocument()
      expect(
        screen.getByText(`${pageSize + 1} - ${totalCount + 1} artists out of ${totalCount + 1}`)
      ).toBeInTheDocument()
    }
  )

  it('should order the artists', async () => {
    const user = userEvent.setup()

    const initialOrder = artistsOrders[1]
    const newOrder = artistsOrders[0]

    let orderBy: string[]
    server.use(
      http.get('/artists', (req) => {
        orderBy = new URL(req.request.url).searchParams.getAll('orderBy')
        const response: WithTotalCountResponse<Artist> = {
          models: artists,
          totalCount: totalCount
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<Artists />)

    await waitFor(() =>
      expect(orderBy).toStrictEqual([initialOrder.property + ' ' + initialOrder.type])
    )

    await user.click(screen.getByRole('button', { name: 'order-artists' }))
    await user.click(screen.getByRole('button', { name: newOrder.label }))

    await waitFor(() =>
      expect(orderBy).toStrictEqual([
        newOrder.property + ' ' + OrderType.Ascending,
        initialOrder.property + ' ' + initialOrder.type
      ])
    )
  })

  it('should filter the artists by min Songs', async () => {
    const user = userEvent.setup()

    const newMinSongsValue = 1

    let searchBy: string[]
    server.use(
      http.get('/artists', (req) => {
        searchBy = new URL(req.request.url).searchParams.getAll('searchBy')
        const response: WithTotalCountResponse<Artist> = {
          models: artists,
          totalCount: totalCount
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<Artists />)

    await waitFor(() => expect(searchBy).toStrictEqual([]))

    await user.click(screen.getByRole('button', { name: 'filter-artists' }))

    await user.clear(screen.getByRole('textbox', { name: /min songs/i }))
    await user.type(
      screen.getByRole('textbox', { name: /min songs/i }),
      newMinSongsValue.toString()
    )
    await user.click(screen.getByRole('button', { name: 'apply-filters' }))

    await waitFor(() => {
      expect(searchBy).toStrictEqual([
        `${ArtistProperty.Songs} ${FilterOperator.GreaterThanOrEqual} ${newMinSongsValue}`
      ])
      expect(window.location.search).toMatch(/&f=/)
    })
  })
})
