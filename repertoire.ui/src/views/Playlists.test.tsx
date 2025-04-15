import { screen, waitFor } from '@testing-library/react'
import Playlists from './Playlists.tsx'
import { emptyPlaylist, reduxRouterRender } from '../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import Playlist from '../types/models/Playlist.ts'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../state/store.ts'
import createOrder from '../utils/createOrder.ts'
import playlistsOrders from '../data/playlists/playlistsOrders.ts'

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
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => {
    window.location.search = ''
    server.resetHandlers()
  })

  afterAll(() => server.close())

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

    const initialOrder = playlistsOrders[2]
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

    await waitFor(() => expect(orderBy).toEqual([createOrder(initialOrder)]))

    await user.click(screen.getByRole('button', { name: 'order-playlists' }))
    await user.click(screen.getByRole('button', { name: newOrder.label }))

    await waitFor(() => expect(orderBy).toEqual([createOrder(newOrder), createOrder(initialOrder)]))
  })
})
