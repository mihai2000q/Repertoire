import { reduxRender, withToastify } from '../../../test-utils.tsx'
import AddExistingArtistAlbumsModal from './AddExistingArtistAlbumsModal.tsx'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { AddAlbumsToArtistRequest } from '../../../types/requests/ArtistRequests.ts'
import { AlbumSearch } from '../../../types/models/Search.ts'
import SearchType from '../../../types/enums/SearchType.ts'

describe('Add Existing Artist Albums Modal', () => {
  const albums: AlbumSearch[] = [
    {
      id: '1',
      title: 'Album 1',
      imageUrl: 'something.png',
      type: SearchType.Album
    },
    {
      id: '2',
      title: 'Album 2',
      type: SearchType.Album
    },
    {
      id: '3',
      title: 'Album 11',
      type: SearchType.Album
    },
    {
      id: '4',
      title: 'Album 12',
      type: SearchType.Album
    }
  ]

  const handlers = [
    http.get('/search', (req) => {
      const query = new URL(req.request.url).searchParams.get('query')
      let localAlbums = albums
      if (query !== '') {
        localAlbums = localAlbums.filter((album) => album.title.startsWith(query))
      }
      const response: WithTotalCountResponse<AlbumSearch> = {
        models: localAlbums,
        totalCount: localAlbums.length
      }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<AddExistingArtistAlbumsModal opened={true} onClose={() => {}} artistId={''} />)

    expect(screen.getByRole('dialog', { name: /add existing albums/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /add existing albums/i })).toBeInTheDocument()
    expect(screen.getByText(/choose albums/i)).toBeInTheDocument()
    expect(screen.getByLabelText('info-icon')).toBeInTheDocument()
    expect(screen.getByRole('searchbox', { name: /search/i })).toBeInTheDocument()
    expect(screen.getByRole('searchbox', { name: /search/i })).toBeDisabled()
    expect(screen.getByRole('searchbox', { name: /search/i })).toHaveValue('')
    expect(screen.getByRole('button', { name: /add/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /add/i })).toHaveAttribute('data-disabled', 'true')
    expect(screen.getByTestId('albums-loader')).toBeInTheDocument()

    await user.hover(screen.getByRole('button', { name: /add/i }))
    expect(await screen.findByText(/select albums/i)).toBeInTheDocument()

    expect(await screen.findByRole('checkbox', { name: /select all/i })).toBeInTheDocument()
    expect(screen.queryByText(/no albums/i)).not.toBeInTheDocument()
    expect(screen.getByRole('searchbox', { name: /search/i })).not.toBeDisabled()

    expect(screen.getAllByRole('checkbox')).toHaveLength(albums.length + 1) // plus the select all checkbox

    for (const album of albums) {
      const renderedAlbum = screen.getByLabelText(`album-${album.title}`)
      expect(renderedAlbum).toBeInTheDocument()

      expect(screen.getByRole('checkbox', { name: album.title })).toBeInTheDocument()
      expect(screen.getByRole('checkbox', { name: album.title })).not.toBeChecked()
      if (album.imageUrl) {
        expect(screen.getByRole('img', { name: album.title })).toHaveAttribute(
          'src',
          album.imageUrl
        )
      } else {
        expect(screen.getByLabelText(`default-icon-${album.title}`)).toBeInTheDocument()
      }
      expect(screen.getByText(album.title)).toBeInTheDocument()
    }
  })

  it('should show text when there are no albums and hide select all checkbox', async () => {
    server.use(
      http.get('/search', () => {
        const response: WithTotalCountResponse<AlbumSearch> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRender(<AddExistingArtistAlbumsModal opened={true} onClose={() => {}} artistId={''} />)

    expect(await screen.findByText(/no albums/i)).toBeInTheDocument()
    expect(screen.queryByRole('checkbox', { name: /select all/i })).not.toBeInTheDocument()
  })

  it('should send updated query when the search box is filled', async () => {
    const user = userEvent.setup()

    let capturedSearchParams: URLSearchParams
    server.use(
      http.get('/search', (req) => {
        capturedSearchParams = new URL(req.request.url).searchParams
        const response: WithTotalCountResponse<AlbumSearch> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRender(<AddExistingArtistAlbumsModal opened={true} onClose={() => {}} artistId={''} />)

    expect(await screen.findByText(/no albums/i)).toBeInTheDocument()

    expect(capturedSearchParams.get('query')).toBe('')
    expect(capturedSearchParams.get('currentPage')).toBe('1')
    expect(capturedSearchParams.get('pageSize')).toBe('20')
    expect(capturedSearchParams.get('order')).match(/updatedAt:desc/i)
    expect(capturedSearchParams.getAll('filter')).toHaveLength(1)
    expect(capturedSearchParams.getAll('filter')[0]).match(/artist IS NULL/i)

    // search
    const searchValue = 'Album 1'
    await user.type(screen.getByRole('searchbox', { name: /search/i }), searchValue)

    await waitFor(() => {
      expect(capturedSearchParams.get('query')).toBe(searchValue)
    })
  })

  it('should search, select albums and send request when clicking add', async () => {
    const user = userEvent.setup()

    const artistId = 'some-artist-id'
    const albumsToSelect = [albums[0], albums[3]]
    const onClose = vitest.fn()

    let capturedRequest: AddAlbumsToArtistRequest
    server.use(
      http.post('/artists/add-albums', async (req) => {
        capturedRequest = (await req.request.json()) as AddAlbumsToArtistRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <AddExistingArtistAlbumsModal opened={true} onClose={onClose} artistId={artistId} />
      )
    )

    // it happens that the other album intended to select has a similar name
    await user.type(
      await screen.findByRole('searchbox', { name: /search/i }),
      albumsToSelect[0].title
    )

    for (const album of albumsToSelect) {
      await user.click(await screen.findByRole('checkbox', { name: album.title }))
    }

    const addButton = screen.getByRole('button', { name: /add/i })
    expect(addButton).not.toHaveAttribute('data-disabled')
    await user.click(addButton)

    expect(await screen.findByText(/albums added to artist/i)).toBeInTheDocument()
    expect(screen.getByRole('searchbox', { name: /search/i })).toHaveValue('')

    expect(capturedRequest).toStrictEqual({
      id: artistId,
      albumIds: albumsToSelect.map((album) => album.id)
    })
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should select all albums and then deselect them', async () => {
    const user = userEvent.setup()

    reduxRender(<AddExistingArtistAlbumsModal opened={true} onClose={() => {}} artistId={''} />)

    await user.click(await screen.findByRole('checkbox', { name: /select all/i }))
    screen.getAllByRole('checkbox').forEach((c) => expect(c).toBeChecked())

    expect(screen.queryByRole('checkbox', { name: /select all/i })).not.toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: /deselect all/i })).toBeInTheDocument()

    await user.click(await screen.findByRole('checkbox', { name: /deselect all/i }))
    screen.getAllByRole('checkbox').forEach((c) => expect(c).not.toBeChecked())
  })

  it('should search for albums based on title', async () => {
    const user = userEvent.setup()

    reduxRender(<AddExistingArtistAlbumsModal opened={true} onClose={() => {}} artistId={''} />)

    expect(await screen.findAllByLabelText(/album-/i)).toHaveLength(albums.length)

    const searchBox = screen.getByRole('searchbox', { name: /search/i })

    await user.type(searchBox, 'Album 2')

    await waitFor(() => {
      expect(screen.getByTestId('loading-overlay-fetching')).toBeVisible()
    })
    await waitFor(() => {
      expect(screen.queryByTestId('loading-overlay-fetching')).not.toBeInTheDocument()
    })

    expect(await screen.findAllByLabelText(/album-/i)).toHaveLength(1)

    await user.clear(searchBox)
    await user.type(searchBox, 'gibberish')

    expect(await screen.findByText(/no albums/i)).toBeInTheDocument()
    expect(screen.queryByRole('checkbox', { name: /select all/i })).not.toBeInTheDocument()
    expect(screen.queryAllByLabelText(/album-/i)).toHaveLength(0)
  })

  it('should deselect albums that do not match the search value, if they were previously selected', async () => {
    const user = userEvent.setup()

    const albumToNotDeselect = albums[0].title
    const albumToDeselect = albums[1].title

    reduxRender(<AddExistingArtistAlbumsModal opened={true} onClose={() => {}} artistId={''} />)

    const searchBox = screen.getByRole('searchbox', { name: /search/i })

    // check albums
    await user.click(await screen.findByRole('checkbox', { name: albumToNotDeselect }))
    await user.click(await screen.findByRole('checkbox', { name: albumToDeselect }))

    // search for the first album, so that second one disappears
    await user.type(searchBox, albumToNotDeselect)

    await waitFor(() => {
      expect(screen.getByTestId('loading-overlay-fetching')).toBeVisible()
    })
    await waitFor(() => {
      expect(screen.queryByTestId('loading-overlay-fetching')).not.toBeInTheDocument()
    })

    expect(screen.queryByRole('checkbox', { name: albumToDeselect })).not.toBeInTheDocument()

    // clear the search so that the second album re-appears (unchecked)
    await user.clear(searchBox)

    expect(await screen.findByRole('checkbox', { name: albumToDeselect })).not.toBeChecked()
    expect(await screen.findByRole('checkbox', { name: albumToNotDeselect })).toBeChecked()
  })
})
