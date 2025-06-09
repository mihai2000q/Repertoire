import {
  emptyAlbum,
  emptyArtist,
  emptySong,
  reduxRender,
  withToastify
} from '../../../test-utils.tsx'
import AddPlaylistSongsModal from './AddPlaylistSongsModal.tsx'
import Song from '../../../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { screen, waitFor, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { AddSongsToPlaylistRequest } from '../../../types/requests/PlaylistRequests.ts'
import Artist from '../../../types/models/Artist.ts'
import Album from '../../../types/models/Album.ts'
import { expect } from 'vitest'
import OrderType from '../../../types/enums/OrderType.ts'
import SongProperty from '../../../types/enums/SongProperty.ts'
import FilterOperator from '../../../types/enums/FilterOperator.ts'

describe('Add Playlist Songs Modal', () => {
  const album: Album = {
    ...emptyAlbum,
    id: '1',
    title: 'Album'
  }

  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist'
  }

  const songs: Song[] = [
    {
      ...emptySong,
      id: '1',
      title: 'Song 1',
      imageUrl: 'something.png',
      artist: artist
    },
    {
      ...emptySong,
      id: '2',
      title: 'Song 2',
      imageUrl: 'something.png',
      album: {
        ...album,
        imageUrl: 'something-album.png'
      },
      artist: artist
    },
    {
      ...emptySong,
      id: '3',
      title: 'Song 11',
      album: {
        ...album,
        imageUrl: 'something-album.png'
      }
    },
    {
      ...emptySong,
      id: '4',
      title: 'Song 12',
      imageUrl: 'something.png'
    },
    {
      ...emptySong,
      id: '5',
      title: 'Song 512',
      imageUrl: 'something.png'
    },
    {
      ...emptySong,
      id: '6',
      title: 'Song 6'
    }
  ]

  const handlers = [
    http.get('/songs', (req) => {
      const searchBy = new URL(req.request.url).searchParams.getAll('searchBy')
      let localSongs = songs
      if (searchBy.length === 2) {
        const searchValue = searchBy[1].replace('songs.title ~* ', '').replaceAll("'", '')
        localSongs = localSongs.filter((song) => song.title.startsWith(searchValue))
      }
      const response: WithTotalCountResponse<Song> = {
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

    reduxRender(<AddPlaylistSongsModal opened={true} onClose={() => {}} playlistId={''} />)

    expect(screen.getByRole('dialog', { name: /add playlist songs/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /add playlist songs/i })).toBeInTheDocument()
    expect(screen.getByText(/choose songs/i)).toBeInTheDocument()
    expect(screen.getByRole('searchbox', { name: /search/i })).toBeInTheDocument()
    expect(screen.getByRole('searchbox', { name: /search/i })).toBeDisabled()
    expect(screen.getByRole('searchbox', { name: /search/i })).toHaveValue('')
    expect(screen.getByRole('button', { name: /add/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /add/i })).toHaveAttribute('data-disabled', 'true')
    expect(screen.getByTestId('songs-loader')).toBeInTheDocument()

    await user.hover(screen.getByRole('button', { name: /add/i }))
    expect(await screen.findByText(/select songs/i)).toBeInTheDocument()

    expect(await screen.findByRole('checkbox', { name: 'Select all' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /show all/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /show all/i })).not.toBeChecked()
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
      if (song.album) expect(within(renderedSong).getByText(song.album.title)).toBeInTheDocument()
      if (song.artist) expect(within(renderedSong).getByText(song.artist.name)).toBeInTheDocument()
    }
  })

  it('should show text when there are no songs and hide select all checkbox', async () => {
    server.use(
      http.get('/songs', () => {
        const response: WithTotalCountResponse<Song> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRender(<AddPlaylistSongsModal opened={true} onClose={() => {}} playlistId={''} />)

    expect(await screen.findByText(/no songs/i)).toBeInTheDocument()
    expect(screen.queryByRole('checkbox', { name: 'Select all' })).not.toBeInTheDocument()
  })

  it('should send updated query when the search box is filled', async () => {
    const user = userEvent.setup()

    const playlistId = 'some-playlist-id'

    let capturedSearchParams: URLSearchParams
    server.use(
      http.get('/songs', (req) => {
        capturedSearchParams = new URL(req.request.url).searchParams
        const response: WithTotalCountResponse<Song> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRender(<AddPlaylistSongsModal opened={true} onClose={() => {}} playlistId={playlistId} />)

    expect(await screen.findByText(/no songs/i)).toBeInTheDocument()

    expect(capturedSearchParams.get('currentPage')).toBe('1')
    expect(capturedSearchParams.get('pageSize')).toBe('20')
    expect(capturedSearchParams.get('orderBy')).toBe(
      `${SongProperty.LastModified} ${OrderType.Descending}`
    )
    expect(capturedSearchParams.getAll('searchBy')).toHaveLength(1)
    expect(capturedSearchParams.getAll('searchBy')[0]).toBe(
      `${SongProperty.PlaylistId} ${FilterOperator.NotEqualVariant} ${playlistId}`
    )

    // search
    const searchValue = 'Song 1'
    await user.type(screen.getByRole('searchbox', { name: /search/i }), searchValue)

    await waitFor(() => {
      expect(capturedSearchParams.getAll('searchBy')).toHaveLength(2)
      expect(capturedSearchParams.getAll('searchBy')[0]).toBe(
        `${SongProperty.PlaylistId} ${FilterOperator.NotEqualVariant} ${playlistId}`
      )
      expect(capturedSearchParams.getAll('searchBy')[1]).toBe(
        `${SongProperty.Title} ${FilterOperator.PatternMatching} ${searchValue}`
      )
    })
  })

  it('should search, select songs and send request when clicking add', async () => {
    const user = userEvent.setup()

    const playlistId = 'some-playlist-id'
    const songsToSelect = [songs[0], songs[3]]
    const onClose = vitest.fn()

    let capturedRequest: AddSongsToPlaylistRequest
    server.use(
      http.post('/playlists/add-songs', async (req) => {
        capturedRequest = (await req.request.json()) as AddSongsToPlaylistRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <AddPlaylistSongsModal opened={true} onClose={onClose} playlistId={playlistId} />
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

    expect(await screen.findByText(/songs added to playlist/i)).toBeInTheDocument()
    expect(screen.getByRole('searchbox', { name: /search/i })).toHaveValue('')

    expect(capturedRequest).toStrictEqual({
      id: playlistId,
      songIds: songsToSelect.map((song) => song.id)
    })
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should select all songs and then deselect them', async () => {
    const user = userEvent.setup()

    reduxRender(<AddPlaylistSongsModal opened={true} onClose={() => {}} playlistId={''} />)

    await user.click(await screen.findByRole('checkbox', { name: 'Select all' }))
    screen.getAllByRole('checkbox').forEach((c) => expect(c).toBeChecked())

    expect(screen.queryByRole('checkbox', { name: 'Select all' })).not.toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: 'Deselect all' })).toBeInTheDocument()

    await user.click(await screen.findByRole('checkbox', { name: 'Deselect all' }))
    screen.getAllByRole('checkbox').forEach((c) => expect(c).not.toBeChecked())
  })

  it('should search for songs based on title', async () => {
    const user = userEvent.setup()

    reduxRender(<AddPlaylistSongsModal opened={true} onClose={() => {}} playlistId={''} />)

    expect(await screen.findAllByLabelText(/song-/i)).toHaveLength(songs.length)

    const searchBox = screen.getByRole('searchbox', { name: /search/i })

    await user.type(searchBox, 'Song 2')

    await waitFor(() => expect(screen.getAllByLabelText(/song-/i)).toHaveLength(1))

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

    reduxRender(<AddPlaylistSongsModal opened={true} onClose={() => {}} playlistId={''} />)

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

  it('should send songs request with filter by playlist by default and without when showing all', async () => {
    const user = userEvent.setup()

    const playlistId = 'some-id'

    let capturedSearchParams: URLSearchParams
    server.use(
      http.get('/songs', (req) => {
        capturedSearchParams = new URL(req.request.url).searchParams
        const response: WithTotalCountResponse<Song> = {
          models: songs,
          totalCount: songs.length
        }
        return HttpResponse.json(response)
      })
    )

    reduxRender(<AddPlaylistSongsModal opened={true} onClose={() => {}} playlistId={playlistId} />)

    await waitFor(() => {
      expect(capturedSearchParams.getAll('searchBy')).toHaveLength(1)
      expect(capturedSearchParams.getAll('searchBy')[0]).toBe(
        `${SongProperty.PlaylistId} ${FilterOperator.NotEqualVariant} ${playlistId}`
      )
    })

    await user.click(await screen.findByRole('button', { name: /show all/i }))
    expect(screen.getByRole('button', { name: /show all/i })).toBeChecked()

    await waitFor(() => expect(capturedSearchParams.getAll('searchBy')).toHaveLength(0))
  })
})
