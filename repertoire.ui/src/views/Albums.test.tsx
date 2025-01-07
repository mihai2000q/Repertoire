import { screen } from '@testing-library/react'
import Albums from './Albums.tsx'
import { reduxRouterRender } from '../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { userEvent } from '@testing-library/user-event'
import Artist from '../types/models/Artist.ts'
import Album from '../types/models/Album.ts'
import Song from '../types/models/Song.ts'

describe('Albums', () => {
  const albums: Album[] = [
    {
      id: '1',
      title: 'Album 1',
      songs: [],
      createdAt: '',
      updatedAt: ''
    },
    {
      id: '2',
      title: 'Album 2',
      songs: [],
      createdAt: '',
      updatedAt: ''
    }
  ]
  const totalCount = 2

  const initialCurrentPage = 1
  const pageSize = 20

  const getSongsWithoutAlbums = () =>
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

  const getAlbumsWithPagination = (totalCount: number = 30) =>
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
    http.get('/artists', async () => {
      const response: WithTotalCountResponse<Artist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display relevant info when there are albums', async () => {
    const [{ container }] = reduxRouterRender(<Albums />)

    expect(screen.getByRole('heading', { name: /albums/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /new-album/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order-albums/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /filter-albums/i })).toBeInTheDocument()
    expect(screen.getByTestId('albums-loader')).toBeInTheDocument()
    expect(container.querySelector('.mantine-Loader-root')).toBeInTheDocument()

    expect(await screen.findByLabelText('new-album-card')).toBeInTheDocument()
    expect(screen.getAllByLabelText(/album-card-/)).toHaveLength(albums.length)
    expect(
      screen.getByText(`${initialCurrentPage} - ${albums.length} albums out of ${totalCount}`)
    ).toBeInTheDocument()
    albums.forEach((album) =>
      expect(screen.getByLabelText(`album-card-${album.title}`)).toBeInTheDocument()
    )
    expect(screen.getByTestId('albums-pagination')).toBeInTheDocument()
  })

  it('should open the add new album modal when clicking the new album button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRouterRender(<Albums />)

    const newAlbumButton = screen.getByRole('button', { name: /new-album/i })
    await user.click(newAlbumButton)

    // Assert
    expect(await screen.findByRole('heading', { name: /add new album/i })).toBeInTheDocument()
  })

  it('should open the add new album modal when clicking the new album card button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRouterRender(<Albums />)

    const newAlbumCardButton = await screen.findByLabelText('new-album-card')
    await user.click(newAlbumCardButton)

    // Assert
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
    expect(screen.queryByLabelText('new-album-card')).toBeInTheDocument()
    expect(screen.queryAllByLabelText(/album-card-/)).toHaveLength(0)
    expect(screen.queryByTestId('albums-pagination')).toBeInTheDocument()
  })

  it('should paginate the albums', async () => {
    // Arrange
    const user = userEvent.setup()

    const totalCount = 30

    server.use(getAlbumsWithPagination(totalCount))

    // Act
    reduxRouterRender(<Albums />)

    // Assert
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
  })

  it('the new album card should not be displayed on first page, but on the last', async () => {
    // Arrange
    const user = userEvent.setup()

    server.use(getAlbumsWithPagination())

    // Act
    reduxRouterRender(<Albums />)

    // Assert
    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('new-album-card')).not.toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('new-album-card')).toBeInTheDocument()
  })

  it('should display unknown album card when there are songs without album', async () => {
    // Arrange
    server.use(getSongsWithoutAlbums())

    // Act
    reduxRouterRender(<Albums />)

    // Assert
    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('unknown-album-card')).toBeInTheDocument()
    expect(
      screen.getByText(
        `${initialCurrentPage} - ${albums.length + 1} albums out of ${totalCount + 1}`
      )
    ).toBeInTheDocument()
  })

  it('should display unknown album card when there are songs without album on the last page, but not on the first', async () => {
    // Arrange
    const user = userEvent.setup()

    const totalCount = 30

    server.use(getSongsWithoutAlbums(), getAlbumsWithPagination(totalCount))

    // Act
    reduxRouterRender(<Albums />)

    // Assert
    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('unknown-album-card')).not.toBeInTheDocument()
    expect(
      screen.getByText(`${initialCurrentPage} - ${pageSize} albums out of ${totalCount + 1}`)
    ).toBeInTheDocument()

    const pagination = screen.getByRole('button', { name: '2' })
    await user.click(pagination)

    expect(await screen.findByTestId('albums-pagination')).toBeInTheDocument()
    expect(screen.queryByLabelText('unknown-album-card')).toBeInTheDocument()
    expect(
      screen.getByText(`${pageSize + 1} - ${totalCount + 1} albums out of ${totalCount + 1}`)
    ).toBeInTheDocument()
  })
})