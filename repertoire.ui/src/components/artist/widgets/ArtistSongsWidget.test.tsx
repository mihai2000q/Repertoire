import { emptyAlbum, emptySong, reduxRouterRender } from '../../../test-utils.tsx'
import ArtistSongsWidget from './ArtistSongsWidget.tsx'
import Song from '../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import Order from 'src/types/Order.ts'
import artistSongsOrders from '../../../data/artist/artistSongsOrders.ts'
import { SongSearch } from '../../../types/models/Search.ts'
import { createRef } from 'react'
import { beforeEach } from 'vitest'

describe('Artist Songs Widget', () => {
  const songModels: Song[] = [
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

  const songs: WithTotalCountResponse<Song> = {
    models: songModels,
    totalCount: songModels.length
  }

  const artistId = '1'

  const order: Order = {
    label: "Songs' order",
    property: 'songs_order'
  }

  const handlers = [
    http.get('/search', async () => {
      const response: WithTotalCountResponse<SongSearch> = {
        models: [],
        totalCount: 0
      }
      return HttpResponse.json(response)
    }),
    http.get(`/songs/guitar-tunings`, () => {
      return HttpResponse.json([])
    }),
    http.get(`/playlists`, () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  beforeEach(() => vi.clearAllMocks())

  afterEach(() => {
    vi.restoreAllMocks()
    server.resetHandlers()
  })

  beforeAll(() => {
    server.listen()
    // Mock Context
    vi.mock('../../../context/MainContext.tsx', () => ({
      useMain: vi.fn(() => ({
        ref: createRef()
      }))
    }))
  })

  afterAll(() => server.close())

  it('should render and display songs', () => {
    reduxRouterRender(
      <ArtistSongsWidget
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    expect(screen.queryByTestId('songs-loader')).not.toBeInTheDocument()
    expect(screen.getByLabelText('songs-widget')).toBeInTheDocument()
    expect(screen.getByText('Songs')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: order.label })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'songs-more-menu' })).toBeInTheDocument()
    expect(screen.getAllByLabelText(/song-card-/)).toHaveLength(songs.models.length)
    songs.models.forEach((song) =>
      expect(screen.getByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    )
    expect(screen.getByLabelText('new-songs-widget')).toBeInTheDocument()
  })

  it('should display loader on loading', async () => {
    reduxRouterRender(
      <ArtistSongsWidget
        songs={songs}
        artistId={artistId}
        isLoading={true}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    expect(screen.getByTestId('songs-loader')).toBeInTheDocument()
    expect(screen.queryByLabelText('songs-widget')).not.toBeInTheDocument()
  })

  it('should display orders and be able to change it', async () => {
    const user = userEvent.setup()

    const newOrder = artistSongsOrders[0]
    const setOrder = vitest.fn()

    reduxRouterRender(
      <ArtistSongsWidget
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={setOrder}
        isUnknownArtist={false}
      />
    )

    await user.click(screen.getByRole('button', { name: order.label }))

    artistSongsOrders.forEach((o) =>
      expect(screen.getByRole('menuitem', { name: o.label })).toBeInTheDocument()
    )

    await user.click(screen.getByRole('menuitem', { name: newOrder.label }))

    expect(setOrder).toHaveBeenCalledWith(newOrder)
  })

  it('should display more menu', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistSongsWidget
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))

    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add existing songs/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should display less information on the more menu, when the artist is unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistSongsWidget
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={true}
      />
    )

    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))

    expect(screen.queryByRole('menuitem', { name: /perfect rehearsal/i })).not.toBeInTheDocument()
    expect(screen.queryByRole('menuitem', { name: /add existing songs/i })).not.toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new song/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open add existing songs modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistSongsWidget
          songs={songs}
          artistId={artistId}
          isLoading={false}
          order={order}
          setOrder={() => {}}
          isUnknownArtist={false}
        />
      )

      await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /add existing songs/i }))

      expect(await screen.findByRole('dialog', { name: /add existing songs/i })).toBeInTheDocument()
    })

    it('should open add new song modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistSongsWidget
          songs={songs}
          artistId={artistId}
          isLoading={false}
          order={order}
          setOrder={() => {}}
          isUnknownArtist={false}
        />
      )

      await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /add new song/i }))

      expect(await screen.findByRole('dialog', { name: /add new song/i })).toBeInTheDocument()
    })
  })

  it('should open Add existing songs modal, when clicking new song card and the artist is not unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistSongsWidget
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    await user.click(screen.getByLabelText('new-songs-widget'))

    expect(await screen.findByRole('dialog', { name: /add existing songs/i })).toBeInTheDocument()
  })

  it('should open Add new song modal, when clicking new song card and the artist is unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistSongsWidget
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={true}
      />
    )

    await user.click(screen.getByLabelText('new-songs-widget'))

    expect(await screen.findByRole('dialog', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should show the drawer when selecting songs, and the context menu when right-clicking after selection', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistSongsWidget
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={true}
      />
    )

    // selection drawer
    await user.keyboard('{Control>}')
    await user.click(screen.getByLabelText(`song-card-${songModels[0].title}`))
    await user.keyboard('{/Control}')
    expect(screen.getByLabelText('songs-selection-drawer')).toBeInTheDocument()

    // context menu
    // await user.pointer({ keys: '[MouseRight>]', target: screen.getByLabelText(`song-card-${songModels[0].title}`) })
    // expect(await screen.findByRole('menu', { name: 'songs-context-menu' })).toBeInTheDocument()
  })
})
