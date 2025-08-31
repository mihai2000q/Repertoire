import { screen, waitFor } from '@testing-library/react'
import Albums from './Albums.tsx'
import {
  defaultAlbumFiltersMetadata,
  emptyAlbum,
  emptySong,
  reduxRouterRender
} from '../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { userEvent } from '@testing-library/user-event'
import Album from '../types/models/Album.ts'
import Song from '../types/models/Song.ts'
import { RootState } from '../state/store.ts'
import { SearchBase } from '../types/models/Search.ts'
import albumsOrders from '../data/albums/albumsOrders.ts'
import FilterOperator from '../types/enums/FilterOperator.ts'
import AlbumProperty from '../types/enums/properties/AlbumProperty.ts'
import { expect } from 'vitest'
import OrderType from '../types/enums/OrderType.ts'
import Playlist from '../types/models/Playlist.ts'

describe('Albums', () => {
  const albums: Album[] = [
    {
      ...emptyAlbum,
      id: '1',
      title: 'Album 1'
    },
    {
      ...emptyAlbum,
      id: '2',
      title: 'Album 2'
    }
  ]
  const totalCount = 2

  const initialCurrentPage = 1
  const pageSize = 40

  const getSongsWithoutAlbums = () =>
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

  const getAlbumsWithPagination = (totalCount: number = 50) =>
    http.get('/albums', async (req) => {
      const currentPage = new URL(req.request.url).searchParams.get('currentPage')
      const response: WithTotalCountResponse<Album> =
        parseInt(currentPage) === 1
          ? {
              models: Array.from({ length: pageSize }).map((_, i) => ({
                ...albums[0],
                id: i.toString()
              })),
              totalCount: totalCount
            }
          : {
              models: Array.from({ length: totalCount - pageSize }).map((_, i) => ({
                ...albums[0],
                id: i.toString()
              })),
              totalCount: totalCount
            }
      return HttpResponse.json(response)
    })

  const getNoAlbums = () =>
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    })

  const handlers = [
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = {
        models: albums,
        totalCount: totalCount
      }
      return HttpResponse.json(response)
    }),
    http.get('/search', async () => {
      const response: WithTotalCountResponse<SearchBase> = { models: [], totalCount: 0 }
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
    http.get('/albums/filters-metadata', async () => {
      return HttpResponse.json(defaultAlbumFiltersMetadata)
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

  it('should render and display relevant info when there are albums', async () => {
    const [_, store] = reduxRouterRender(<Albums />)

    expect((store.getState() as RootState).global.documentTitle).toMatch(/albums/i)
    expect(screen.getByRole('heading', { name: /albums/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /new-album/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-albums/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-albums/i })).toBeDisabled()
    expect(screen.getByRole('button', { name: /filter-albums/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /filter-albums/i })).toBeDisabled()

    expect(await screen.findByLabelText('new-album-card')).toBeInTheDocument()
    expect(screen.getAllByLabelText(/album-card-/)).toHaveLength(albums.length)
    expect(
      screen.getByText(`${initialCurrentPage} - ${albums.length} albums out of ${totalCount}`)
    ).toBeInTheDocument()
    albums.forEach((album) =>
      expect(screen.getByLabelText(`album-card-${album.title}`)).toBeInTheDocument()
    )
    expect(screen.getByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-albums/i })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: /filter-albums/i })).not.toBeDisabled()
  })

  it('should open the add new album modal when clicking the new album button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<Albums />)

    const newAlbumButton = screen.getByRole('button', { name: /new-album/i })
    await user.click(newAlbumButton)

    expect(await screen.findByRole('heading', { name: /add new album/i })).toBeInTheDocument()
  })

  it('should open the add new album modal when clicking the new album card button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<Albums />)

    const newAlbumCardButton = await screen.findByLabelText('new-album-card')
    await user.click(newAlbumCardButton)

    expect(await screen.findByRole('heading', { name: /add new album/i })).toBeInTheDocument()
  })

  it('should display not display some info when there are no albums', async () => {
    reduxRouterRender(<Albums />)

    server.use(getNoAlbums())

    expect(await screen.findByText(/no albums/)).toBeInTheDocument()
    expect(screen.queryByLabelText('new-album-card')).not.toBeInTheDocument()
    expect(screen.queryByLabelText('unknown-album-card')).not.toBeInTheDocument()
    expect(screen.queryAllByLabelText(/album-card-/)).toHaveLength(0)
    expect(screen.queryByTestId('albums-pagination')).not.toBeInTheDocument()
  })

  it('should display pagination and new album card when the unknown album card is shown and there are no albums', async () => {
    reduxRouterRender(<Albums />)

    server.use(getNoAlbums(), getSongsWithoutAlbums())

    expect(await screen.findByLabelText('unknown-album-card')).toBeInTheDocument()
    expect(screen.queryByText(/no albums/)).not.toBeInTheDocument()
    expect(screen.getByLabelText('new-album-card')).toBeInTheDocument()
    expect(screen.queryAllByLabelText(/album-card-/)).toHaveLength(0)
    expect(screen.getByTestId('albums-pagination')).toBeInTheDocument()
  })

  it('should paginate the albums', async () => {
    const user = userEvent.setup()

    const totalCount = 50

    server.use(getAlbumsWithPagination(totalCount))

    reduxRouterRender(<Albums />)

    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.queryAllByLabelText(/album-card-/)).toHaveLength(pageSize)
    expect(
      screen.getByText(`${initialCurrentPage} - ${pageSize} albums out of ${totalCount}`)
    ).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '1' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '2' })).toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.queryAllByLabelText(/album-card-/)).toHaveLength(totalCount - pageSize)
    expect(
      screen.getByText(`${pageSize + 1} - ${totalCount} albums out of ${totalCount}`)
    ).toBeInTheDocument()
    expect(window.location.search).toContain('p=2')
  })

  it('should display unknown album card when there are songs without album', async () => {
    server.use(getSongsWithoutAlbums())

    reduxRouterRender(<Albums />)

    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('unknown-album-card')).toBeInTheDocument()
    expect(
      screen.getByText(
        `${initialCurrentPage} - ${albums.length + 1} albums out of ${totalCount + 1}`
      )
    ).toBeInTheDocument()
  })

  it('should display the new album card on first page, but on the last', async () => {
    const user = userEvent.setup()

    server.use(getAlbumsWithPagination())

    reduxRouterRender(<Albums />)

    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('new-album-card')).not.toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.getByLabelText('new-album-card')).toBeInTheDocument()
  })

  it('should display the new album and unknown album cards on the last page, but not on the first page', async () => {
    const user = userEvent.setup()

    const totalCount = 50

    server.use(getAlbumsWithPagination(totalCount), getSongsWithoutAlbums())

    reduxRouterRender(<Albums />)

    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('unknown-album-card')).not.toBeInTheDocument()
    expect(screen.queryByLabelText('new-album-card')).not.toBeInTheDocument()
    expect(
      screen.getByText(`${initialCurrentPage} - ${pageSize} albums out of ${totalCount + 1}`)
    ).toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.getByLabelText('unknown-album-card')).toBeInTheDocument()
    expect(screen.getByLabelText('new-album-card')).toBeInTheDocument()
    expect(
      screen.getByText(`${pageSize + 1} - ${totalCount + 1} albums out of ${totalCount + 1}`)
    ).toBeInTheDocument()
  })

  it('should order the albums', async () => {
    const user = userEvent.setup()

    const initialOrder = albumsOrders[1]
    const newOrder = albumsOrders[0]

    let orderBy: string[]
    server.use(
      http.get('/albums', (req) => {
        orderBy = new URL(req.request.url).searchParams.getAll('orderBy')
        const response: WithTotalCountResponse<Album> = {
          models: albums,
          totalCount: totalCount
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<Albums />)

    await waitFor(() =>
      expect(orderBy).toStrictEqual([initialOrder.property + ' ' + initialOrder.type])
    )

    await user.click(screen.getByRole('button', { name: 'order-albums' }))
    await user.click(screen.getByRole('button', { name: newOrder.label }))

    await waitFor(() =>
      expect(orderBy).toStrictEqual([
        newOrder.property + ' ' + OrderType.Ascending,
        initialOrder.property + ' ' + initialOrder.type
      ])
    )
  })

  it('should filter the albums by min Songs', async () => {
    const user = userEvent.setup()

    const newMinSongsValue = 1

    let searchBy: string[]
    server.use(
      http.get('/albums', (req) => {
        searchBy = new URL(req.request.url).searchParams.getAll('searchBy')
        const response: WithTotalCountResponse<Album> = {
          models: albums,
          totalCount: totalCount
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<Albums />)

    await waitFor(() => expect(searchBy).toStrictEqual([]))

    await user.click(screen.getByRole('button', { name: 'filter-albums' }))

    await user.clear(screen.getByRole('textbox', { name: /min songs/i }))
    await user.type(
      screen.getByRole('textbox', { name: /min songs/i }),
      newMinSongsValue.toString()
    )
    await user.click(screen.getByRole('button', { name: 'apply-filters' }))

    await waitFor(() => {
      expect(searchBy).toStrictEqual([
        `${AlbumProperty.Songs} ${FilterOperator.GreaterThanOrEqual} ${newMinSongsValue}`
      ])
      expect(window.location.search).toMatch(/&f=/)
    })
  })
})
