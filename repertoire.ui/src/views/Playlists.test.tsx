import { screen } from '@testing-library/react'
import Playlists from './Playlists.tsx'
import { reduxRouterRender } from '../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import Playlist from '../types/models/Playlist.ts'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { userEvent } from '@testing-library/user-event'

describe('Playlists', () => {
  const playlists: Playlist[] = [
    {
      id: '1',
      title: 'Playlist 1',
      description: '',
      songs: [],
      createdAt: '',
      updatedAt: ''
    },
    {
      id: '2',
      title: 'Playlist 2',
      description: '',
      songs: [],
      createdAt: '',
      updatedAt: ''
    }
  ]
  const totalCount = 2

  const initialCurrentPage = 1
  const pageSize = 20

  const getSongsWithPagination = (totalCount: number = 30) =>
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
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display relevant info when there are playlists', async () => {
    const [{ container }] = reduxRouterRender(<Playlists />)

    expect(screen.getByRole('heading', { name: /playlists/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /new-playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-playlists/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /filter-playlists/i })).toBeInTheDocument()
    expect(screen.getByTestId('playlists-loader')).toBeInTheDocument()
    expect(container.querySelector('.mantine-Loader-root')).toBeInTheDocument()

    expect(await screen.findByLabelText('new-playlist-card')).toBeInTheDocument()
    expect(screen.getAllByLabelText(/playlist-card-/)).toHaveLength(playlists.length)
    expect(
      screen.getByText(`${initialCurrentPage} - ${playlists.length} playlists out of ${totalCount}`)
    ).toBeInTheDocument()
    playlists.forEach((playlist) =>
      expect(screen.getByLabelText(`playlist-card-${playlist.title}`)).toBeInTheDocument()
    )
    expect(screen.getByTestId('playlists-pagination')).toBeInTheDocument()
  })

  it('should open the add new playlist modal when clicking the new playlist button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRouterRender(<Playlists />)

    const newSongButton = screen.getByRole('button', { name: /new-playlist/i })
    await user.click(newSongButton)

    // Assert
    expect(await screen.findByRole('heading', { name: /add new playlist/i })).toBeInTheDocument()
  })

  it('should open the add new playlist modal when clicking the new playlist card button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRouterRender(<Playlists />)

    const newSongCardButton = await screen.findByLabelText('new-playlist-card')
    await user.click(newSongCardButton)

    // Assert
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
    // Arrange
    const user = userEvent.setup()

    const totalCount = 30

    server.use(getSongsWithPagination(totalCount))

    // Act
    reduxRouterRender(<Playlists />)

    // Assert
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
  })

  it('the new playlist card should not be displayed on first page, but on the last', async () => {
    // Arrange
    const user = userEvent.setup()

    server.use(getSongsWithPagination())

    // Act
    reduxRouterRender(<Playlists />)

    // Assert
    expect(await screen.findByTestId('playlists-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('new-playlist-card')).not.toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('playlists-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('new-playlist-card')).toBeInTheDocument()
  })
})