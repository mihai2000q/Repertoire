import { screen, waitFor } from '@testing-library/react'
import Playlists from './Playlists.tsx'
import { defaultPlaylistFiltersMetadata, emptyPlaylist, reduxRouterRender } from '../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import Playlist from '../types/models/Playlist.ts'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../state/store.ts'
import playlistsOrders from '../data/playlists/playlistsOrders.ts'
import PlaylistProperty from '../types/enums/properties/PlaylistProperty.ts'
import FilterOperator from '../types/enums/FilterOperator.ts'
import OrderType from '../types/enums/OrderType.ts'
import { expect } from 'vitest'
import { createRef } from 'react'

describe('Playlists', () => {
  const playlists: Playlist[] = [
    {
      ...emptyPlaylist,
      id: '1',
      title: 'Playlist 1'
    },
    {
      ...emptyPlaylist,
      id: '2',
      title: 'Playlist 2'
    }
  ]
  const totalCount = 2

  const initialCurrentPage = 1
  const pageSize = 40

  const getSongsWithPagination = (totalCount: number = 50) =>
    http.get('/playlists', async (req) => {
      const currentPage = new URL(req.request.url).searchParams.get('currentPage')
      const response: WithTotalCountResponse<Playlist> =
        parseInt(currentPage) === 1
          ? {
              models: Array.from({ length: pageSize }).map((_, i) => ({
                ...playlists[0],
                id: i.toString()
              })),
              totalCount: totalCount
            }
          : {
              models: Array.from({ length: totalCount - pageSize }).map((_, i) => ({
                ...playlists[0],
                id: i.toString()
              })),
              totalCount: totalCount
            }
      return HttpResponse.json(response)
    })

  const handlers = [
    http.get('/playlists', async () => {
      const response: WithTotalCountResponse<Playlist> = {
        models: playlists,
        totalCount: totalCount
      }
      return HttpResponse.json(response)
    }),
    http.get('/playlists/filters-metadata', async () => {
      return HttpResponse.json(defaultPlaylistFiltersMetadata)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => {
    server.listen()
    // Mock Drag Select
    vi.mock('dragselect', () => ({
      default: vi.fn(
        class {
          start = vi.fn()
          stop = vi.fn()
          getSelection = vi.fn()
          clearSelection = vi.fn()
          setSettings = vi.fn()
          subscribe = vi.fn()
          unsubscribe = vi.fn()
          addSelectables = vi.fn()
          removeSelectables = vi.fn()
        }
      )
    }))
    // Mock Main Context
    vi.mock('../context/MainContext.tsx', () => ({
      useMain: vi.fn(() => ({
        ref: createRef(),
        mainScroll: { ref: createRef() }
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

  it('should render and display relevant info when there are playlists', async () => {
    const [_, store] = reduxRouterRender(<Playlists />)

    expect((store.getState() as RootState).global.documentTitle).toMatch(/playlists/i)
    expect(screen.getByRole('heading', { name: /playlists/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /new-playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-playlists/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-playlists/i })).toBeDisabled()
    expect(screen.getByRole('button', { name: /filter-playlists/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /filter-playlists/i })).toBeDisabled()

    expect(await screen.findByLabelText('new-playlist-card')).toBeInTheDocument()
    expect(screen.getAllByLabelText(/playlist-card-/)).toHaveLength(playlists.length)
    expect(
      screen.getByText(`${initialCurrentPage} - ${playlists.length} playlists out of ${totalCount}`)
    ).toBeInTheDocument()
    playlists.forEach((playlist) =>
      expect(screen.getByLabelText(`playlist-card-${playlist.title}`)).toBeInTheDocument()
    )
    expect(screen.getByTestId('playlists-pagination')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-playlists/i })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: /filter-playlists/i })).not.toBeDisabled()
  })

  it('should open the add new playlist modal when clicking the new playlist button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<Playlists />)

    const newSongButton = screen.getByRole('button', { name: /new-playlist/i })
    await user.click(newSongButton)

    expect(await screen.findByRole('heading', { name: /add new playlist/i })).toBeInTheDocument()
  })

  it('should open the add new playlist modal when clicking the new playlist card button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<Playlists />)

    const newSongCardButton = await screen.findByLabelText('new-playlist-card')
    await user.click(newSongCardButton)

    expect(await screen.findByRole('heading', { name: /add new playlist/i })).toBeInTheDocument()
  })

  it('should display not display some info when there are no playlists', async () => {
    reduxRouterRender(<Playlists />)

    server.use(
      http.get('/playlists', async () => {
        const response: WithTotalCountResponse<Playlist> = { models: [], totalCount: 0 }
        return HttpResponse.json(response)
      })
    )

    expect(await screen.findByText(/no playlists/)).toBeInTheDocument()
    expect(screen.queryByLabelText('new-playlist-card')).not.toBeInTheDocument()
    expect(screen.queryAllByLabelText(/playlist-card-/)).toHaveLength(0)
    expect(screen.queryByTestId('playlists-pagination')).not.toBeInTheDocument()
  })

  it('should paginate the playlists', async () => {
    const user = userEvent.setup()

    const totalCount = 50

    server.use(getSongsWithPagination(totalCount))

    reduxRouterRender(<Playlists />)

    expect(await screen.findByTestId('playlists-pagination')).toBeInTheDocument()
    expect(screen.queryAllByLabelText(/playlist-card-/)).toHaveLength(pageSize)
    expect(
      screen.getByText(`${initialCurrentPage} - ${pageSize} playlists out of ${totalCount}`)
    ).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '1' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '2' })).toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('playlists-pagination')).toBeInTheDocument()
    expect(screen.queryAllByLabelText(/playlist-card-/)).toHaveLength(totalCount - pageSize)
    expect(
      screen.getByText(`${pageSize + 1} - ${totalCount} playlists out of ${totalCount}`)
    ).toBeInTheDocument()
    expect(window.location.search).toContain('p=2')
  })

  it('the new playlist card should not be displayed on first page, but on the last', async () => {
    const user = userEvent.setup()

    server.use(getSongsWithPagination())

    reduxRouterRender(<Playlists />)

    expect(await screen.findByTestId('playlists-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('new-playlist-card')).not.toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('playlists-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('new-playlist-card')).toBeInTheDocument()
  })

  it('should order the playlists', async () => {
    const user = userEvent.setup()

    const initialOrder = playlistsOrders[1]
    const newOrder = playlistsOrders[0]

    let orderBy: string[]
    server.use(
      http.get('/playlists', (req) => {
        orderBy = new URL(req.request.url).searchParams.getAll('orderBy')
        const response: WithTotalCountResponse<Playlist> = {
          models: playlists,
          totalCount: totalCount
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<Playlists />)

    await waitFor(() =>
      expect(orderBy).toStrictEqual([initialOrder.property + ' ' + initialOrder.type])
    )

    await user.click(screen.getByRole('button', { name: 'order-playlists' }))
    await user.click(screen.getByRole('button', { name: newOrder.label }))

    await waitFor(() =>
      expect(orderBy).toStrictEqual([
        newOrder.property + ' ' + OrderType.Ascending,
        initialOrder.property + ' ' + initialOrder.type
      ])
    )
  })

  it('should filter the playlists by min Songs', async () => {
    const user = userEvent.setup()

    const newMinSongsValue = 1

    let searchBy: string[]
    server.use(
      http.get('/playlists', (req) => {
        searchBy = new URL(req.request.url).searchParams.getAll('searchBy')
        const response: WithTotalCountResponse<Playlist> = {
          models: playlists,
          totalCount: totalCount
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<Playlists />)

    await waitFor(() => expect(searchBy).toStrictEqual([]))

    await user.click(screen.getByRole('button', { name: 'filter-playlists' }))

    await user.clear(screen.getByRole('textbox', { name: /min songs/i }))
    await user.type(
      screen.getByRole('textbox', { name: /min songs/i }),
      newMinSongsValue.toString()
    )
    await user.click(screen.getByRole('button', { name: 'apply-filters' }))

    await waitFor(() => {
      expect(searchBy).toStrictEqual([
        `${PlaylistProperty.Songs} ${FilterOperator.GreaterThanOrEqual} ${newMinSongsValue}`
      ])
      expect(window.location.search).toMatch(/&f=/)
    })
  })

  it.skip('should show the drawer when selecting playlists, and the context menu when right-clicking after selection', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<Playlists />)

    const playlistsArea = screen.getByTestId('playlists-area')
    await user.pointer([
      { keys: '[MouseLeft>]', target: playlistsArea, coords: { x: 0, y: 0 } },
      { coords: { x: 1000, y: 300 } },
      { keys: '[/MouseLeft]' }
    ])

    expect(
      await screen.findByRole('dialog', { name: 'playlists-selection-drawer' })
    ).toBeInTheDocument()

    await user.pointer({ keys: '[MouseRight>]', target: playlistsArea })
    expect(await screen.findByRole('menu', { name: 'playlists-context-menu' })).toBeInTheDocument()
  })
})
