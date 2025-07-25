import {
  defaultSongFiltersMetadata,
  emptyAlbum,
  emptySong,
  reduxRouterRender
} from '../../test-utils.tsx'
import AlbumSongsCard from './AlbumSongsCard.tsx'
import Song from '../../types/models/Song.ts'
import Album from '../../types/models/Album.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import albumSongsOrders from '../../data/album/albumSongsOrders.ts'
import { expect } from 'vitest'
import { SongSearch } from '../../types/models/Search.ts'

describe('Album Songs Card', () => {
  const songs: Song[] = [
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
    }
  ]

  const album: Album = {
    ...emptyAlbum,
    id: '1',
    title: 'Song 1',
    songs: [
      ...songs,
      {
        ...emptySong,
        id: '5',
        title: 'Song 5'
      }
    ]
  }

  const order = albumSongsOrders[1]

  const handlers = [
    http.get('/search', async () => {
      const response: WithTotalCountResponse<SongSearch> = {
        models: [],
        totalCount: 0
      }
      return HttpResponse.json(response)
    }),
    http.get('/playlists', async () => {
      return HttpResponse.json(defaultSongFiltersMetadata)
    }),
    http.get('/songs/guitar-tunings', async () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it("should render and display album's songs when the album is not unknown", async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <AlbumSongsCard
        album={album}
        songs={[]}
        isUnknownAlbum={false}
        order={order}
        setOrder={() => {}}
      />
    )

    expect(screen.getByText(/songs/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'songs-more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: order.label })).toBeInTheDocument()
    expect(screen.getAllByLabelText(/song-card-/)).toHaveLength(album.songs.length)
    album.songs.forEach((song) =>
      expect(screen.getByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    )
    expect(screen.queryByLabelText('new-song-card')).not.toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: order.label }))
    albumSongsOrders.forEach((o) => {
      const screenOrder = screen.getByRole('menuitem', { name: o.label })
      expect(screenOrder).toBeInTheDocument()
      expect(screenOrder).not.toBeDisabled()
    })

    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))
    expect(screen.getByRole('menuitem', { name: /add existing songs/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should render and display the songs when the album is unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <AlbumSongsCard
        album={album}
        songs={songs}
        isUnknownAlbum={true}
        order={order}
        setOrder={() => {}}
      />
    )

    expect(screen.getByText(/songs/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'songs-more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: order.label })).toBeInTheDocument()
    expect(screen.getAllByLabelText(/song-card-/)).toHaveLength(songs.length)
    songs.forEach((song) =>
      expect(screen.getByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    )
    expect(screen.getByLabelText('new-song-card')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: order.label }))
    albumSongsOrders.forEach((o) => {
      const screenOrder = screen.getByRole('menuitem', { name: o.label })
      expect(screenOrder).toBeInTheDocument()
      if (o === albumSongsOrders[0]) {
        expect(screenOrder).toBeDisabled()
      } else {
        expect(screenOrder).not.toBeDisabled()
      }
    })

    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))
    expect(screen.queryByRole('menuitem', { name: /add existing songs/i })).not.toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should display orders and be able to change it', async () => {
    const user = userEvent.setup()

    const newOrder = albumSongsOrders[3]
    const setOrder = vitest.fn()

    reduxRouterRender(
      <AlbumSongsCard
        album={album}
        songs={songs}
        isUnknownAlbum={true}
        order={order}
        setOrder={setOrder}
      />
    )

    await user.click(screen.getByRole('button', { name: order.label }))

    albumSongsOrders.forEach((o) =>
      expect(screen.getByRole('menuitem', { name: o.label })).toBeInTheDocument()
    )

    await user.click(screen.getByRole('menuitem', { name: newOrder.label }))

    expect(setOrder).toHaveBeenCalledWith(newOrder)
  })

  describe('on menu', () => {
    it('should open add existing songs modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <AlbumSongsCard
          album={album}
          songs={[]}
          isUnknownAlbum={false}
          order={order}
          setOrder={() => {}}
        />
      )

      await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /add existing songs/i }))

      expect(await screen.findByRole('dialog', { name: /add existing songs/i })).toBeInTheDocument()
    })

    it('should open add new song modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <AlbumSongsCard
          album={album}
          songs={[]}
          isUnknownAlbum={false}
          order={order}
          setOrder={() => {}}
        />
      )

      await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /add new song/i }))

      expect(await screen.findByRole('dialog', { name: /add new song/i })).toBeInTheDocument()
    })
  })

  it('should display new song card when there are no album songs and open Add existing songs modal', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <AlbumSongsCard
        album={{ ...album, songs: [] }}
        songs={[]}
        isUnknownAlbum={false}
        order={order}
        setOrder={() => {}}
      />
    )

    expect(screen.getByLabelText('new-song-card')).toBeInTheDocument()

    await user.click(screen.getByLabelText('new-song-card'))

    expect(await screen.findByRole('dialog', { name: /add existing songs/i })).toBeInTheDocument()
  })

  it('should display new song card when the album is unknown and open Add new song modal', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <AlbumSongsCard
        album={album}
        songs={songs}
        isUnknownAlbum={true}
        order={order}
        setOrder={() => {}}
      />
    )

    expect(screen.getByLabelText('new-song-card')).toBeInTheDocument()

    await user.click(screen.getByLabelText('new-song-card'))

    expect(await screen.findByRole('dialog', { name: /add new song/i })).toBeInTheDocument()
  })

  it.skip('should be able to reorder songs', () => {})
})
