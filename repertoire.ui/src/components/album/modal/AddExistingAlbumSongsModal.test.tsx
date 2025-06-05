import { reduxRender, withToastify } from '../../../test-utils.tsx'
import AddExistingAlbumSongsModal from './AddExistingAlbumSongsModal.tsx'
import Song from '../../../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { screen, waitFor, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { AddSongsToAlbumRequest } from '../../../types/requests/AlbumRequests.ts'
import SearchType from '../../../types/enums/SearchType.ts'
import { SongSearch } from '../../../types/models/Search.ts'

describe('Add Existing Album Songs Modal', () => {
  const songs: SongSearch[] = [
    {
      id: '1',
      title: 'Song 1',
      imageUrl: 'something.png',
      album: {
        id: '1',
        title: 'Album 1',
        imageUrl: 'something-album.png'
      },
      type: SearchType.Song
    },
    {
      id: '2',
      title: 'Song 2',
      album: {
        id: '2',
        title: 'Album 2',
        imageUrl: 'something-album.png'
      },
      artist: {
        id: '1',
        name: 'Artist'
      },
      type: SearchType.Song
    },
    {
      id: '3',
      title: 'Song 11',
      imageUrl: 'something.png',
      type: SearchType.Song
    },
    {
      id: '4',
      title: 'Song 12',
      type: SearchType.Song
    }
  ]

  const handlers = [
    http.get('/search', (req) => {
      const query = new URL(req.request.url).searchParams.get('query')
      let localSongs = songs
      if (query !== '') {
        localSongs = localSongs.filter((song) => song.title.startsWith(query))
      }
      const response: WithTotalCountResponse<SongSearch> = {
        models: localSongs,
        totalCount: localSongs.length
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

    reduxRender(
      <AddExistingAlbumSongsModal opened={true} onClose={() => {}} albumId={''} artistId={''} />
    )

    expect(screen.getByRole('dialog', { name: /add existing songs/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /add existing songs/i })).toBeInTheDocument()
    expect(screen.getByText(/choose songs/i)).toBeInTheDocument()
    expect(screen.getByLabelText('info-icon')).toBeInTheDocument()
    expect(screen.getByRole('searchbox', { name: /search/i })).toBeInTheDocument()
    expect(screen.getByRole('searchbox', { name: /search/i })).toBeDisabled()
    expect(screen.getByRole('searchbox', { name: /search/i })).toHaveValue('')
    expect(screen.getByRole('button', { name: /add/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /add/i })).toHaveAttribute('data-disabled', 'true')
    expect(screen.getByTestId('songs-loader')).toBeInTheDocument()

    await user.hover(screen.getByRole('button', { name: /add/i }))
    expect(await screen.findByText(/select songs/i)).toBeInTheDocument()

    expect(await screen.findByRole('checkbox', { name: 'Select all' })).toBeInTheDocument()
    expect(screen.queryByText(/no songs/i)).not.toBeInTheDocument()
    expect(screen.getByRole('searchbox', { name: /search/i })).not.toBeDisabled()

    expect(screen.getAllByRole('checkbox')).toHaveLength(songs.length + 1) // plus the select all checkbox

    for (const song of songs) {
      const renderedSong = screen.getByLabelText(`song-${song.title}`)
      expect(renderedSong).toBeInTheDocument()

      expect(screen.getByRole('checkbox', { name: song.title })).toBeInTheDocument()
      expect(screen.getByRole('checkbox', { name: song.title })).not.toBeChecked()
      if (song.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', song.imageUrl)
      } else if (song.album?.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute(
          'src',
          song.album.imageUrl
        )
      } else {
        expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
      }
      expect(screen.getByText(song.title)).toBeInTheDocument()
      if (song.artist) expect(within(renderedSong).getByText(song.artist.name)).toBeInTheDocument()
    }
  })

  it('should show text when there are no songs and hide select all checkbox', async () => {
    server.use(
      http.get('/search', () => {
        const response: WithTotalCountResponse<SongSearch> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRender(
      <AddExistingAlbumSongsModal opened={true} onClose={() => {}} albumId={''} artistId={''} />
    )

    expect(await screen.findByText(/no songs/i)).toBeInTheDocument()
    expect(screen.queryByRole('checkbox', { name: 'Select all' })).not.toBeInTheDocument()
  })

  it('should send updated query when the search box is filled', async () => {
    const user = userEvent.setup()

    let capturedSearchParams: URLSearchParams
    server.use(
      http.get('/search', (req) => {
        capturedSearchParams = new URL(req.request.url).searchParams
        const response: WithTotalCountResponse<Song> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    const [{ rerender }] = reduxRender(
      <AddExistingAlbumSongsModal opened={true} onClose={() => {}} albumId={''} artistId={''} />
    )

    expect(await screen.findByText(/no songs/i)).toBeInTheDocument()

    expect(capturedSearchParams.get('query')).toBe('')
    expect(capturedSearchParams.get('currentPage')).toBe('1')
    expect(capturedSearchParams.get('pageSize')).toBe('20')
    expect(capturedSearchParams.get('order')).match(/updatedAt:desc/i)
    expect(capturedSearchParams.get('type')).toBe(SearchType.Song)
    expect(capturedSearchParams.getAll('filter')).toHaveLength(2)
    expect(capturedSearchParams.getAll('filter')[0]).match(/album IS NULL/i)
    expect(capturedSearchParams.getAll('filter')[1]).match(/artist IS NULL/i)

    // rerender with artist
    const artistId = 'some-artist-id'

    rerender(
      <AddExistingAlbumSongsModal
        opened={true}
        onClose={() => {}}
        albumId={''}
        artistId={artistId}
      />
    )

    await waitFor(() => {
      expect(capturedSearchParams.getAll('filter')[1]).match(
        new RegExp(`artist IS NULL OR artist.id = ${artistId}`, 'i')
      )
    })

    // search
    const searchValue = 'Song 1'
    await user.type(screen.getByRole('searchbox', { name: /search/i }), searchValue)

    await waitFor(() => {
      expect(capturedSearchParams.get('query')).toBe(searchValue)
    })
  })

  it('should search, select songs and send request when clicking add', async () => {
    const user = userEvent.setup()

    const albumId = 'some-album-id'
    const songsToSelect = [songs[0], songs[3]]
    const onClose = vitest.fn()

    let capturedRequest: AddSongsToAlbumRequest
    server.use(
      http.post('/albums/add-songs', async (req) => {
        capturedRequest = (await req.request.json()) as AddSongsToAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <AddExistingAlbumSongsModal
          opened={true}
          onClose={onClose}
          albumId={albumId}
          artistId={''}
        />
      )
    )

    // it happens that the other song intended to select has a similar name
    await user.type(
      await screen.findByRole('searchbox', { name: /search/i }),
      songsToSelect[0].title
    )

    for (const song of songsToSelect.reverse()) {
      await user.click(await screen.findByRole('checkbox', { name: song.title }))
    }

    const addButton = screen.getByRole('button', { name: /add/i })
    expect(addButton).not.toHaveAttribute('data-disabled')
    await user.click(addButton)

    expect(await screen.findByText(/songs added to album/i)).toBeInTheDocument()
    expect(screen.getByRole('searchbox', { name: /search/i })).toHaveValue('')

    expect(capturedRequest).toStrictEqual({
      id: albumId,
      songIds: songsToSelect.map((song) => song.id)
    })
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should select all songs and then deselect them', async () => {
    const user = userEvent.setup()

    reduxRender(
      <AddExistingAlbumSongsModal opened={true} onClose={() => {}} albumId={''} artistId={''} />
    )

    await user.click(await screen.findByRole('checkbox', { name: 'Select all' }))
    screen.getAllByRole('checkbox').forEach((c) => expect(c).toBeChecked())

    expect(screen.queryByRole('checkbox', { name: 'Select all' })).not.toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: 'Deselect all' })).toBeInTheDocument()

    await user.click(await screen.findByRole('checkbox', { name: 'Deselect all' }))
    screen.getAllByRole('checkbox').forEach((c) => expect(c).not.toBeChecked())
  })

  it('should search for songs based on title', async () => {
    const user = userEvent.setup()

    reduxRender(
      <AddExistingAlbumSongsModal opened={true} onClose={() => {}} albumId={''} artistId={''} />
    )

    expect(await screen.findAllByLabelText(/song-/i)).toHaveLength(songs.length)

    const searchBox = screen.getByRole('searchbox', { name: /search/i })

    await user.type(searchBox, 'Song 2')

    await waitFor(() => {
      expect(screen.getByTestId('loading-overlay-fetching')).toBeVisible()
    })
    await waitFor(() => {
      expect(screen.queryByTestId('loading-overlay-fetching')).not.toBeInTheDocument()
    })

    expect(await screen.findAllByLabelText(/song-/i)).toHaveLength(1)

    await user.clear(searchBox)
    await user.type(searchBox, 'gibberish')

    expect(await screen.findByText(/no songs/i)).toBeInTheDocument()
    expect(screen.queryByRole('checkbox', { name: 'Select all' })).not.toBeInTheDocument()
    expect(screen.queryAllByLabelText(/song-/i)).toHaveLength(0)
  })

  it('should deselect songs that do not match the search value, if they were previously selected', async () => {
    const user = userEvent.setup()

    const songToNotDeselect = songs[0].title
    const songToDeselect = songs[1].title

    reduxRender(
      <AddExistingAlbumSongsModal opened={true} onClose={() => {}} albumId={''} artistId={''} />
    )

    const searchBox = screen.getByRole('searchbox', { name: /search/i })

    // check songs
    await user.click(await screen.findByRole('checkbox', { name: songToNotDeselect }))
    await user.click(await screen.findByRole('checkbox', { name: songToDeselect }))

    // search for the first song, so that the second one disappears
    await user.type(searchBox, songToNotDeselect)

    await waitFor(() => {
      expect(screen.getByTestId('loading-overlay-fetching')).toBeVisible()
    })
    await waitFor(() => {
      expect(screen.queryByTestId('loading-overlay-fetching')).not.toBeInTheDocument()
    })

    expect(screen.queryByRole('checkbox', { name: songToDeselect })).not.toBeInTheDocument()

    // clear the search so that the second song re-appears (unchecked)
    await user.clear(searchBox)

    expect(await screen.findByRole('checkbox', { name: songToDeselect })).not.toBeChecked()
    expect(await screen.findByRole('checkbox', { name: songToNotDeselect })).toBeChecked()
  })
})
