import { emptyAlbum, emptyOrder, emptyPlaylist, emptySong, reduxRouterRender } from '../../test-utils.tsx'
import PlaylistSongsCard from './PlaylistSongsCard.tsx'
import Song from '../../types/models/Song.ts'
import Playlist from '../../types/models/Playlist.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'

describe('Playlist Songs Card', () => {
  const playlist: Playlist = {
    ...emptyPlaylist,
    id: '1',
    title: 'Song 1',
    songs: [
      {
        ...emptySong,
        id: '1',
        title: 'Song 1',
        imageUrl: 'something.png',
        album: {
          ...emptyAlbum,
          imageUrl: 'something-album.png'
        }
      },
      {
        ...emptySong,
        id: '2',
        title: 'Song 2',
        album: {
          ...emptyAlbum,
          imageUrl: 'something-album.png'
        }
      },
      {
        ...emptySong,
        id: '3',
        title: 'Song 3',
        imageUrl: 'something.png',
        album: {
          ...emptyAlbum
        }
      },
      {
        ...emptySong,
        id: '4',
        title: 'Song 4',
        album: {
          ...emptyAlbum
        }
      },
      {
        ...emptySong,
        id: '5',
        title: 'Song 5',
        imageUrl: 'something.png'
      },
      {
        ...emptySong,
        id: '6',
        title: 'Song 6'
      }
    ]
  }

  const handlers = [
    http.get('/songs', async () => {
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

  it("should render and display playlist's songs", () => {
    reduxRouterRender(
      <PlaylistSongsCard playlist={playlist} order={emptyOrder} setOrder={vi.fn()} />
    )

    expect(screen.getByText(/songs/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'songs-more-menu' })).toBeInTheDocument()
    expect(screen.getAllByLabelText(/song-card-/)).toHaveLength(playlist.songs.length)
    playlist.songs.forEach((song) =>
      expect(screen.getByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    )
    expect(screen.queryByLabelText('new-song-card')).not.toBeInTheDocument()
  })

  it('should display menu', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <PlaylistSongsCard playlist={playlist} order={emptyOrder} setOrder={vi.fn()} />
    )

    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))

    expect(screen.getByRole('menuitem', { name: /add songs/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open add playlist songs modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <PlaylistSongsCard playlist={playlist} order={emptyOrder} setOrder={vi.fn()} />
      )

      await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /add songs/i }))

      expect(await screen.findByRole('dialog', { name: /add playlist songs/i })).toBeInTheDocument()
    })
  })

  it('should display new song card when there are no playlist songs and open Add playlist songs modal', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <PlaylistSongsCard
        playlist={{ ...playlist, songs: [] }}
        order={emptyOrder}
        setOrder={vi.fn()}
      />
    )

    expect(screen.getByLabelText('new-song-card')).toBeInTheDocument()

    await user.click(screen.getByLabelText('new-song-card'))

    expect(await screen.findByRole('dialog', { name: /add playlist songs/i })).toBeInTheDocument()
  })

  it.skip('should be able to reorder', () => {})
})
