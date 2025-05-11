import { emptyAlbum, emptyOrder, reduxRouterRender } from '../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Album from 'src/types/models/Album.ts'
import ArtistAlbumCard from './ArtistAlbumCard.tsx'
import dayjs from 'dayjs'
import { expect } from 'vitest'
import AlbumProperty from '../../types/enums/AlbumProperty.ts'
import { RemoveAlbumsFromArtistRequest } from '../../types/requests/ArtistRequests.ts'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'

describe('Artist Album Card', () => {
  const album: Album = {
    ...emptyAlbum,
    id: '1',
    title: 'Album 1'
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => {
    server.resetHandlers()
    window.location.pathname = '/'
  })

  afterAll(() => server.close())

  it('should render and display minimal information', async () => {
    reduxRouterRender(
      <ArtistAlbumCard album={album} artistId={''} isUnknownArtist={false} order={emptyOrder} />
    )

    expect(screen.getByLabelText(`default-icon-${album.title}`)).toBeInTheDocument()
    expect(screen.getByText(album.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should render and display maximal information', async () => {
    const localAlbum: Album = {
      ...album,
      imageUrl: 'something.png'
    }

    reduxRouterRender(
      <ArtistAlbumCard
        album={localAlbum}
        artistId={''}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    expect(screen.getByRole('img', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.title })).toHaveAttribute(
      'src',
      localAlbum.imageUrl
    )
    expect(screen.getByText(localAlbum.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  describe('on order property change', () => {
    it('should display the release date when it is release date', async () => {
      const localAlbum: Album = {
        ...album,
        releaseDate: '2024-10-11'
      }

      reduxRouterRender(
        <ArtistAlbumCard
          album={localAlbum}
          artistId={''}
          isUnknownArtist={false}
          order={{ ...emptyOrder, property: AlbumProperty.ReleaseDate }}
        />
      )

      expect(
        screen.getByText(dayjs(localAlbum.releaseDate).format('D MMM YYYY'))
      ).toBeInTheDocument()
    })

    it('should display the rehearsals, when it is rehearsals', () => {
      const order = {
        ...emptyOrder,
        property: AlbumProperty.Rehearsals
      }

      const localAlbum = {
        ...album,
        rehearsals: 34
      }

      reduxRouterRender(
        <ArtistAlbumCard album={localAlbum} artistId={''} isUnknownArtist={false} order={order} />
      )

      expect(screen.getAllByText(localAlbum.rehearsals)).toHaveLength(2) // the one visible and the one in the tooltip
      expect(screen.getAllByText(localAlbum.rehearsals)[0]).toBeVisible()
      expect(screen.getAllByText(localAlbum.rehearsals)[1]).not.toBeVisible()
    })

    it('should display the confidence bar, when it is confidence', () => {
      const order = {
        ...emptyOrder,
        property: AlbumProperty.Confidence
      }

      const localAlbum = {
        ...album,
        confidence: 23
      }

      reduxRouterRender(
        <ArtistAlbumCard album={localAlbum} artistId={''} isUnknownArtist={false} order={order} />
      )

      expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    })

    it('should display the progress bar, when it is progress', () => {
      const order = {
        ...emptyOrder,
        property: AlbumProperty.Progress
      }

      const localAlbum = {
        ...album,
        progress: 123
      }

      reduxRouterRender(
        <ArtistAlbumCard album={localAlbum} artistId={''} isUnknownArtist={false} order={order} />
      )

      expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    })

    it('should display the last played date, when it is last played', () => {
      const order = {
        ...emptyOrder,
        property: AlbumProperty.LastPlayed
      }

      const localAlbum = {
        ...album,
        lastTimePlayed: '2024-10-12T10:30'
      }

      const [{ rerender }] = reduxRouterRender(
        <ArtistAlbumCard album={localAlbum} artistId={''} isUnknownArtist={false} order={order} />
      )

      expect(
        screen.getByText(dayjs(localAlbum.lastTimePlayed).format('D MMM YYYY'))
      ).toBeInTheDocument()

      rerender(
        <ArtistAlbumCard album={album} artistId={''} isUnknownArtist={false} order={order} />
      )

      expect(screen.getByText(/never/i)).toBeInTheDocument()
    })
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistAlbumCard album={album} artistId={''} isUnknownArtist={false} order={emptyOrder} />
    )

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`album-card-${album.title}`)
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove from artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistAlbumCard album={album} artistId={''} isUnknownArtist={false} order={emptyOrder} />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove from artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should display less information on the menu when the artist is unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistAlbumCard album={album} artistId={''} isUnknownArtist={true} order={emptyOrder} />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.queryByRole('menuitem', { name: /remove from artist/i })).not.toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to album when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistAlbumCard album={album} artistId={''} isUnknownArtist={false} order={emptyOrder} />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))

      expect(window.location.pathname).toBe(`/album/${album.id}`)
    })

    it('should display warning modal and remove album from artist, when clicking on remove from artist', async () => {
      const user = userEvent.setup()

      let capturedRequest: RemoveAlbumsFromArtistRequest
      server.use(
        http.put('/artists/remove-albums', async (req) => {
          capturedRequest = (await req.request.json()) as RemoveAlbumsFromArtistRequest
          return HttpResponse.json()
        })
      )

      const artistId = 'some-artist-id'

      reduxRouterRender(
        <ArtistAlbumCard
          album={album}
          artistId={artistId}
          isUnknownArtist={false}
          order={emptyOrder}
        />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /remove from artist/i }))

      expect(await screen.findByRole('dialog', { name: /remove album/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /remove album/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(capturedRequest).toStrictEqual({
        id: artistId,
        albumIds: [album.id]
      })
    })

    it('should display warning modal and delete album, when clicking on delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/albums/${album.id}`, () => {
          return HttpResponse.json()
        })
      )

      reduxRouterRender(
        <ArtistAlbumCard album={album} artistId={''} isUnknownArtist={false} order={emptyOrder} />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete album/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete album/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))
    })
  })
})
