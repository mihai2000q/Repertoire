import { emptyAlbum, reduxRouterRender } from '../../../test-utils.tsx'
import ArtistAlbumsCard from './ArtistAlbumsCard.tsx'
import Album from '../../../types/models/Album.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import Order from 'src/types/Order.ts'
import artistAlbumsOrders from '../../../data/artist/artistAlbumsOrders.ts'
import { AlbumSearch } from '../../../types/models/Search.ts'

describe('Artist Albums Card', () => {
  const albumModels: Album[] = [
    {
      ...emptyAlbum,
      id: '1',
      title: 'Album 1'
    },
    {
      ...emptyAlbum,
      id: '2',
      title: 'Album 2'
    }
  ]

  const albums: WithTotalCountResponse<Album> = {
    models: albumModels,
    totalCount: albumModels.length
  }

  const artistId = '1'

  const order: Order = {
    label: "Albums' order",
    property: 'albums_order'
  }

  const handlers = [
    http.get('/search', async () => {
      const response: WithTotalCountResponse<AlbumSearch> = {
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

  it('should render and display albums', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistAlbumsCard
        albums={albums}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    expect(screen.queryByTestId('albums-loader')).not.toBeInTheDocument()
    expect(screen.getByLabelText('albums-card')).toBeInTheDocument()
    expect(screen.getByText('Albums')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: order.label })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'albums-more-menu' })).toBeInTheDocument()
    expect(screen.getAllByLabelText(/album-card-/)).toHaveLength(albums.models.length)
    albums.models.forEach((album) =>
      expect(screen.getByLabelText(`album-card-${album.title}`)).toBeInTheDocument()
    )
    expect(screen.getByLabelText('new-albums-card')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'albums-more-menu' }))

    expect(screen.getByRole('menuitem', { name: /add existing albums/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new album/i })).toBeInTheDocument()
  })

  it('should display loader on loading', async () => {
    reduxRouterRender(
      <ArtistAlbumsCard
        albums={albums}
        artistId={artistId}
        isLoading={true}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    expect(screen.getByTestId('albums-loader')).toBeInTheDocument()
    expect(screen.queryByLabelText('albums-card')).not.toBeInTheDocument()
  })

  it('should display orders and be able to change it', async () => {
    const user = userEvent.setup()

    const newOrder = artistAlbumsOrders[0]
    const setOrder = vitest.fn()

    reduxRouterRender(
      <ArtistAlbumsCard
        albums={albums}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={setOrder}
        isUnknownArtist={false}
      />
    )

    await user.click(screen.getByRole('button', { name: order.label }))

    artistAlbumsOrders.forEach((o) =>
      expect(screen.getByRole('menuitem', { name: o.label })).toBeInTheDocument()
    )

    await user.click(screen.getByRole('menuitem', { name: newOrder.label }))

    expect(setOrder).toHaveBeenCalledWith(newOrder)
  })

  it('should display more menu', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistAlbumsCard
        albums={albums}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    await user.click(screen.getByRole('button', { name: 'albums-more-menu' }))

    expect(screen.getByRole('menuitem', { name: /add existing albums/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new album/i })).toBeInTheDocument()
  })

  it('should display less information on the more menu, when the artist is unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistAlbumsCard
        albums={albums}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={true}
      />
    )

    await user.click(screen.getByRole('button', { name: 'albums-more-menu' }))

    expect(screen.queryByRole('menuitem', { name: /add existing albums/i })).not.toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new album/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open add existing albums modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistAlbumsCard
          albums={albums}
          artistId={artistId}
          isLoading={false}
          order={order}
          setOrder={() => {}}
          isUnknownArtist={false}
        />
      )

      await user.click(screen.getByRole('button', { name: 'albums-more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /add existing albums/i }))

      expect(
        await screen.findByRole('dialog', { name: /add existing albums/i })
      ).toBeInTheDocument()
    })

    it('should open add new album modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistAlbumsCard
          albums={albums}
          artistId={artistId}
          isLoading={false}
          order={order}
          setOrder={() => {}}
          isUnknownArtist={false}
        />
      )

      await user.click(screen.getByRole('button', { name: 'albums-more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /add new album/i }))

      expect(await screen.findByRole('dialog', { name: /add new album/i })).toBeInTheDocument()
    })
  })

  it('should open Add existing albums modal, when clicking new album card and the artist is not unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistAlbumsCard
        albums={albums}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={false}
      />
    )

    await user.click(screen.getByLabelText('new-albums-card'))

    expect(await screen.findByRole('dialog', { name: /add existing albums/i })).toBeInTheDocument()
  })

  it('should open Add new album modal, when clicking new album card and the artist is unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistAlbumsCard
        albums={albums}
        artistId={artistId}
        isLoading={false}
        order={order}
        setOrder={() => {}}
        isUnknownArtist={true}
      />
    )

    await user.click(screen.getByLabelText('new-albums-card'))

    expect(await screen.findByRole('dialog', { name: /add new album/i })).toBeInTheDocument()
  })
})
