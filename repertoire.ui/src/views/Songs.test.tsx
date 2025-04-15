import { screen, waitFor } from '@testing-library/react'
import Songs from './Songs.tsx'
import { emptySong, reduxRouterRender } from '../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import Song, { GuitarTuning, SongSectionType } from '../types/models/Song.ts'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../state/store.ts'
import { SearchBase } from '../types/models/Search.ts'
import createOrder from '../utils/createOrder.ts'
import songsOrders from '../data/songs/songsOrders.ts'

describe('Songs', () => {
  const songs: Song[] = [
    {
      ...emptySong,
      id: '1',
      title: 'All for justice'
    },
    {
      ...emptySong,
      id: '2',
      title: 'Seek and Destroy'
    }
  ]
  const totalCount = 2

  const initialCurrentPage = 1
  const pageSize = 40

  const getSongsWithPagination = (totalCount: number = 50) =>
    http.get('/songs', async (req) => {
      const currentPage = new URL(req.request.url).searchParams.get('currentPage')
      const response: WithTotalCountResponse<Song> =
        parseInt(currentPage) === 1
          ? {
              models: Array.from({ length: pageSize }).map((_, i) => ({
                ...songs[0],
                id: i.toString()
              })),
              totalCount: totalCount
            }
          : {
              models: Array.from({ length: totalCount - pageSize }).map((_, i) => ({
                ...songs[0],
                id: i.toString()
              })),
              totalCount: totalCount
            }
      return HttpResponse.json(response)
    })

  const handlers = [
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
        models: songs,
        totalCount: totalCount
      }
      return HttpResponse.json(response)
    }),
    http.get('/search', async () => {
      const response: WithTotalCountResponse<SearchBase> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/songs/guitar-tunings', async () => {
      const response: GuitarTuning[] = []
      return HttpResponse.json(response)
    }),
    http.get('/songs/sections/types', async () => {
      const response: SongSectionType[] = []
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

  it('should render and display relevant info when there are songs', async () => {
    const [_, store] = reduxRouterRender(<Songs />)

    expect((store.getState() as RootState).global.documentTitle).toMatch(/songs/i)
    expect(screen.getByRole('heading', { name: /songs/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /new-song/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-songs/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-songs/i })).toBeDisabled()
    expect(screen.getByRole('button', { name: /filter-songs/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /filter-songs/i })).toBeDisabled()

    expect(await screen.findByLabelText('new-song-card')).toBeInTheDocument()
    expect(screen.getAllByLabelText(/song-card-/)).toHaveLength(songs.length)
    expect(
      screen.getByText(`${initialCurrentPage} - ${songs.length} songs out of ${totalCount}`)
    ).toBeInTheDocument()
    songs.forEach((song) =>
      expect(screen.getByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    )
    expect(screen.getByTestId('songs-pagination')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-songs/i })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: /filter-songs/i })).not.toBeDisabled()
  })

  it('should open the add new song modal when clicking the new song button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<Songs />)

    const newSongButton = screen.getByRole('button', { name: /new-song/i })
    await user.click(newSongButton)

    expect(await screen.findByRole('heading', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should open the add new song modal when clicking the new song card button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<Songs />)

    const newSongCardButton = await screen.findByLabelText('new-song-card')
    await user.click(newSongCardButton)

    expect(await screen.findByRole('heading', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should display not display some info when there are no songs', async () => {
    reduxRouterRender(<Songs />)

    server.use(
      http.get('/songs', async () => {
        const response: WithTotalCountResponse<Song> = { models: [], totalCount: 0 }
        return HttpResponse.json(response)
      })
    )

    expect(await screen.findByText(/no songs/)).toBeInTheDocument()
    expect(screen.queryByLabelText('new-song-card')).not.toBeInTheDocument()
    expect(screen.queryAllByLabelText(/song-card-/)).toHaveLength(0)
    expect(screen.queryByTestId('songs-pagination')).not.toBeInTheDocument()
  })

  it('should paginate the songs', async () => {
    const user = userEvent.setup()

    const totalCount = 50

    server.use(getSongsWithPagination(totalCount))

    reduxRouterRender(<Songs />)

    expect(await screen.findByTestId('songs-pagination')).toBeInTheDocument()
    expect(screen.queryAllByLabelText(/song-card-/)).toHaveLength(pageSize)
    expect(
      screen.getByText(`${initialCurrentPage} - ${pageSize} songs out of ${totalCount}`)
    ).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '1' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '2' })).toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('songs-pagination')).toBeInTheDocument()
    expect(screen.queryAllByLabelText(/song-card-/)).toHaveLength(totalCount - pageSize)
    expect(
      screen.getByText(`${pageSize + 1} - ${totalCount} songs out of ${totalCount}`)
    ).toBeInTheDocument()
    expect(window.location.search).toContain('p=2')
  })

  it('the new song card should not be displayed on first page, but on the last', async () => {
    const user = userEvent.setup()

    server.use(getSongsWithPagination())

    reduxRouterRender(<Songs />)

    expect(await screen.findByTestId('songs-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('new-song-card')).not.toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('songs-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('new-song-card')).toBeInTheDocument()
  })

  it('should order the songs', async () => {
    const user = userEvent.setup()

    const initialOrder = songsOrders[8]
    const newOrder = songsOrders[0]

    let orderBy: string[]
    server.use(
      http.get('/songs', (req) => {
        orderBy = new URL(req.request.url).searchParams.getAll('orderBy')
        const response: WithTotalCountResponse<Song> = {
          models: songs,
          totalCount: totalCount
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<Songs />)

    await waitFor(() => expect(orderBy).toEqual([createOrder(initialOrder)]))

    await user.click(screen.getByRole('button', { name: 'order-songs' }))
    await user.click(screen.getByRole('button', { name: newOrder.label }))

    await waitFor(() => expect(orderBy).toEqual([createOrder(newOrder), createOrder(initialOrder)]))
  })
})
