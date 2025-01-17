import { screen } from '@testing-library/react'
import Artists from './Artists.tsx'
import { reduxRouterRender } from '../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { userEvent } from '@testing-library/user-event'
import Artist from '../types/models/Artist.ts'
import Song from '../types/models/Song.ts'
import Album from '../types/models/Album.ts'

describe('Artists', () => {
  const artists: Artist[] = [
    {
      id: '1',
      name: 'Artist 1',
      albums: [],
      songs: [],
      createdAt: '',
      updatedAt: ''
    },
    {
      id: '2',
      name: 'Artist 2',
      albums: [],
      songs: [],
      createdAt: '',
      updatedAt: ''
    }
  ]
  const totalCount = 2

  const initialCurrentPage = 1
  const pageSize = 20

  const getSongsWithoutArtists = () =>
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
        models: [
          {
            id: '',
            title: 'Some song',
            description: '',
            isRecorded: false,
            sections: [],
            rehearsals: 0,
            confidence: 0,
            progress: 0,
            createdAt: '',
            updatedAt: ''
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
            id: '',
            title: 'Some Album',
            songs: [],
            createdAt: '',
            updatedAt: ''
          }
        ],
        totalCount: 1
      }
      return HttpResponse.json(response)
    })

  const getArtistsWithPagination = (totalCount: number = 30) =>
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
    http.get('/artists', async () => {
      const response: WithTotalCountResponse<Artist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display relevant info when there are artists', async () => {
    reduxRouterRender(<Artists />)

    expect(screen.getByRole('heading', { name: /artists/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /new-artist/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-artists/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /filter-artists/i })).toBeInTheDocument()
    expect(screen.getByTestId('artists-loader')).toBeInTheDocument()

    expect(await screen.findByLabelText('new-artist-card')).toBeInTheDocument()
    expect(screen.getAllByLabelText(/artist-card-/)).toHaveLength(artists.length)
    expect(
      screen.getByText(`${initialCurrentPage} - ${artists.length} artists out of ${totalCount}`)
    ).toBeInTheDocument()
    artists.forEach((artist) =>
      expect(screen.getByLabelText(`artist-card-${artist.name}`)).toBeInTheDocument()
    )
    expect(screen.getByTestId('artists-pagination')).toBeInTheDocument()
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

    const totalCount = 30

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
  })

  it('the new artist card should not be displayed on first page, but on the last', async () => {
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
    'should display unknown artist card when there are songs or albums without artist on the last page, but not on the first',
    async (withAlbum, withSong) => {
      const user = userEvent.setup()

      const totalCount = 30

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
      expect(
        screen.getByText(`${initialCurrentPage} - ${pageSize} artists out of ${totalCount + 1}`)
      ).toBeInTheDocument()

      const pagination = screen.getByRole('button', { name: '2' })
      await user.click(pagination)

      expect(await screen.findByTestId('artists-pagination')).toBeInTheDocument()
      expect(screen.queryByLabelText('unknown-artist-card')).toBeInTheDocument()
      expect(
        screen.getByText(`${pageSize + 1} - ${totalCount + 1} artists out of ${totalCount + 1}`)
      ).toBeInTheDocument()
    }
  )
})
