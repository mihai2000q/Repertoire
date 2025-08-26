import {
  emptyAlbum,
  emptyPlaylist,
  emptySong,
  reduxRouterRender,
  withToastify
} from '../../test-utils.tsx'
import PlaylistSongsWidget from './PlaylistSongsWidget.tsx'
import Song from '../../types/models/Song.ts'
import Playlist from '../../types/models/Playlist.ts'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { MainScrollProvider } from '../../context/MainScrollContext.tsx'
import playlistSongsOrders from '../../data/playlist/playlistSongsOrders.ts'
import OrderType from '../../types/enums/OrderType.ts'
import { ShufflePlaylistSongsRequest } from '../../types/requests/PlaylistRequests.ts'

describe('Playlist Songs Card', () => {
  const songs: Song[] = [
    {
      ...emptySong,
      playlistSongId: '1',
      title: 'Song 1',
      imageUrl: 'something.png',
      album: {
        ...emptyAlbum,
        imageUrl: 'something-album.png'
      }
    },
    {
      ...emptySong,
      playlistSongId: '2',
      title: 'Song 2',
      album: {
        ...emptyAlbum,
        imageUrl: 'something-album.png'
      }
    },
    {
      ...emptySong,
      playlistSongId: '3',
      title: 'Song 3',
      imageUrl: 'something.png',
      album: {
        ...emptyAlbum
      }
    },
    {
      ...emptySong,
      playlistSongId: '4',
      title: 'Song 4',
      album: {
        ...emptyAlbum
      }
    },
    {
      ...emptySong,
      playlistSongId: '5',
      title: 'Song 5',
      imageUrl: 'something.png'
    },
    {
      ...emptySong,
      playlistSongId: '6',
      title: 'Song 6'
    }
  ]

  const playlist: Playlist = {
    ...emptyPlaylist,
    id: '1',
    title: 'Song 1'
  }

  const handlers = [
    http.get(`/playlist/${playlist.id}`, async () => {
      return HttpResponse.json(playlist)
    }),
    http.get(`/playlists/songs/${playlist.id}`, async () => {
      const response: WithTotalCountResponse<Song> = {
        models: songs,
        totalCount: songs.length
      }
      return HttpResponse.json(response)
    }),
    http.get(`/songs`, async () => {
      const response: WithTotalCountResponse<Song> = {
        models: [],
        totalCount: 0
      }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const render = () =>
    reduxRouterRender(
      withToastify(
        <MainScrollProvider>
          <PlaylistSongsWidget playlistId={playlist.id} />
        </MainScrollProvider>
      )
    )

  it("should render and display playlist's songs", async () => {
    render()

    expect(await screen.findByText(/songs/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'songs-more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: playlistSongsOrders[0].label })).toBeInTheDocument()
    expect(screen.getAllByLabelText(/song-card-/)).toHaveLength(songs.length)
    songs.forEach((song) =>
      expect(screen.getByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    )
    expect(screen.queryByLabelText('new-song-card')).not.toBeInTheDocument()
  })

  it('should display orders and be able to change it', async () => {
    const user = userEvent.setup()

    const newOrder = playlistSongsOrders[1]

    let capturedSearchParams: URLSearchParams
    server.use(
      http.get(`/playlists/songs/${playlist.id}`, async (req) => {
        capturedSearchParams = new URL(req.request.url).searchParams
        const response: WithTotalCountResponse<Song> = {
          models: songs,
          totalCount: songs.length
        }
        return HttpResponse.json(response)
      })
    )

    render()

    await user.click(await screen.findByRole('button', { name: playlistSongsOrders[0].label }))

    playlistSongsOrders.forEach((o) =>
      expect(screen.getByRole('menuitem', { name: o.label })).toBeInTheDocument()
    )

    await user.click(screen.getByRole('menuitem', { name: newOrder.label }))

    await waitFor(() => {
      expect(capturedSearchParams.getAll('orderBy')).toHaveLength(1)
      expect(capturedSearchParams.get('orderBy')).toBe(
        `${newOrder.property} ${newOrder.type ?? OrderType.Ascending}`
      )
    })
  })

  it('should display menu', async () => {
    const user = userEvent.setup()

    render()

    await user.click(await screen.findByRole('button', { name: 'songs-more-menu' }))

    expect(screen.getByRole('menuitem', { name: /shuffle/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add songs/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should send shuffle request', async () => {
      const user = userEvent.setup()

      let capturedRequest: ShufflePlaylistSongsRequest
      server.use(
        http.post(`/playlists/songs/shuffle`, async (req) => {
          capturedRequest = (await req.request.json()) as ShufflePlaylistSongsRequest
          return HttpResponse.json('it worked!')
        })
      )

      render()

      await user.click(await screen.findByRole('button', { name: 'songs-more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /shuffle/i }))
      await user.click(await screen.findByRole('button', { name: /confirm/i }))

      expect(await screen.findByText(/playlist shuffled/i)).toBeInTheDocument()
      expect(capturedRequest.id).toBe(playlist.id)
    })

    it('should open add playlist songs modal', async () => {
      const user = userEvent.setup()

      render()

      await user.click(await screen.findByRole('button', { name: 'songs-more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /add songs/i }))

      expect(await screen.findByRole('dialog', { name: /add playlist songs/i })).toBeInTheDocument()
    })
  })

  it('should display new song card when there are no playlist songs and open Add playlist songs modal', async () => {
    const user = userEvent.setup()

    server.use(
      http.get(`/playlists/songs/${playlist.id}`, async () => {
        const response: WithTotalCountResponse<Song> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    render()

    expect(await screen.findByLabelText('new-song-card')).toBeInTheDocument()

    await user.click(screen.getByLabelText('new-song-card'))

    expect(await screen.findByRole('dialog', { name: /add playlist songs/i })).toBeInTheDocument()
  })

  it.skip('should be able to reorder', () => {})
})
