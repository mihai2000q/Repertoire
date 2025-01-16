import { screen } from '@testing-library/react'
import Songs from './Songs.tsx'
import { reduxRouterRender } from '../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import Song, { GuitarTuning, SongSectionType } from '../types/models/Song.ts'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { userEvent } from '@testing-library/user-event'
import Artist from '../types/models/Artist.ts'
import Album from '../types/models/Album.ts'

describe('Songs', () => {
  const songs: Song[] = [
    {
      id: '1',
      title: 'All for justice',
      description: '',
      isRecorded: false,
      sections: [],
      rehearsals: 0,
      confidence: 0,
      progress: 0,
      createdAt: '',
      updatedAt: ''
    },
    {
      id: '2',
      description: '',
      title: 'Seek and Destroy',
      isRecorded: true,
      sections: [],
      rehearsals: 0,
      confidence: 0,
      progress: 0,
      createdAt: '',
      updatedAt: ''
    }
  ]
  const totalCount = 2

  const initialCurrentPage = 1
  const pageSize = 20

  const getSongsWithPagination = (totalCount: number = 30) =>
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
    http.get('/artists', async () => {
      const response: WithTotalCountResponse<Artist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = { models: [], totalCount: 0 }
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

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display relevant info when there are songs', async () => {
    reduxRouterRender(<Songs />)

    expect(screen.getByRole('heading', { name: /songs/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /new-song/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-songs/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /filter-songs/i })).toBeInTheDocument()
    expect(screen.getByTestId('songs-loader')).toBeInTheDocument()

    expect(await screen.findByLabelText('new-song-card')).toBeInTheDocument()
    expect(screen.getAllByLabelText(/song-card-/)).toHaveLength(songs.length)
    expect(
      screen.getByText(`${initialCurrentPage} - ${songs.length} songs out of ${totalCount}`)
    ).toBeInTheDocument()
    songs.forEach((song) =>
      expect(screen.getByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    )
    expect(screen.getByTestId('songs-pagination')).toBeInTheDocument()
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

    const totalCount = 30

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
})
