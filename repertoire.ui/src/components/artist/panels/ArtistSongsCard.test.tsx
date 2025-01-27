import { emptyAlbum, emptySong, reduxRender } from '../../../test-utils.tsx'
import ArtistSongsCard from './ArtistSongsCard.tsx'
import Song from '../../../types/models/Song.ts'
import { screen, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { RemoveSongsFromArtistRequest } from '../../../types/requests/ArtistRequests.ts'
import Order from 'src/types/Order.ts'
import artistSongsOrders from '../../../data/artist/artistSongsOrders.ts'

describe('Artist Songs Card', () => {
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
    value: 'order value'
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

  it('should render and display songs', async () => {
    const user = userEvent.setup()

    reduxRender(
      <ArtistSongsCard
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    expect(screen.queryByTestId('songs-loader')).not.toBeInTheDocument()
    expect(screen.getByLabelText('songs-card')).toBeInTheDocument()
    expect(screen.getByText('Songs')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: order.label })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'songs-more-menu' })).toBeInTheDocument()
    expect(screen.getAllByLabelText(/song-card-/)).toHaveLength(songs.models.length)
    songs.models.forEach((song) =>
      expect(screen.getByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    )
    expect(screen.getByLabelText('new-songs-card')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))

    expect(screen.getByRole('menuitem', { name: /add existing songs/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should display loader on loading', async () => {
    reduxRender(
      <ArtistSongsCard
        songs={songs}
        artistId={artistId}
        isLoading={true}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    expect(screen.getByTestId('songs-loader')).toBeInTheDocument()
    expect(screen.queryByLabelText('songs-card')).not.toBeInTheDocument()
  })

  it('should display orders and be able to change it', async () => {
    const user = userEvent.setup()

    const newOrder = artistSongsOrders[0]
    const setOrder = vitest.fn()

    reduxRender(
      <ArtistSongsCard
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

    reduxRender(
      <ArtistSongsCard
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))

    expect(screen.getByRole('menuitem', { name: /add existing songs/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should display less information on the more menu, when the artist is unknown', async () => {
    const user = userEvent.setup()

    reduxRender(
      <ArtistSongsCard
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={true}
      />
    )

    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))

    expect(screen.queryByRole('menuitem', { name: /add existing songs/i })).not.toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new song/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open add existing songs modal', async () => {
      const user = userEvent.setup()

      reduxRender(
        <ArtistSongsCard
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

      reduxRender(
        <ArtistSongsCard
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

    reduxRender(
      <ArtistSongsCard
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    await user.click(screen.getByLabelText('new-songs-card'))

    expect(await screen.findByRole('dialog', { name: /add existing songs/i })).toBeInTheDocument()
  })

  it('should open Add new song modal, when clicking new song card and the artist is unknown', async () => {
    const user = userEvent.setup()

    reduxRender(
      <ArtistSongsCard
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={true}
      />
    )

    await user.click(screen.getByLabelText('new-songs-card'))

    expect(await screen.findByRole('dialog', { name: /add new song/i })).toBeInTheDocument()
  })

  it("should send 'remove songs from artist request' when clicking on the more menu of a song card", async () => {
    const user = userEvent.setup()

    const song = songs.models[0]

    let capturedRequest: RemoveSongsFromArtistRequest
    server.use(
      http.put('/artists/remove-songs', async (req) => {
        capturedRequest = (await req.request.json()) as RemoveSongsFromArtistRequest
        return HttpResponse.json()
      })
    )

    reduxRender(
      <ArtistSongsCard
        songs={songs}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    const songCard1 = screen.getByLabelText(`song-card-${song.title}`)

    await user.click(within(songCard1).getByRole('button', { name: 'more-menu' }))
    await user.click(screen.getByRole('menuitem', { name: /remove/i }))
    await user.click(screen.getByRole('button', { name: /yes/i })) // warning modal

    expect(capturedRequest.id).toBe(artistId)
    expect(capturedRequest.songIds).toStrictEqual([song.id])
  })
})
