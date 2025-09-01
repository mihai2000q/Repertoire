import Album from 'src/types/models/Album.ts'
import {
  emptyAlbum,
  emptyArtist,
  reduxRouterRender
} from '../../test-utils.tsx'
import AlbumCard from './AlbumCard.tsx'
import { screen } from '@testing-library/react'
import Artist from 'src/types/models/Artist.ts'
import userEvent from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { RootState } from 'src/state/store.ts'

describe('Album Card', () => {
  const album: Album = {
    ...emptyAlbum,
    id: '1',
    title: 'Album 1'
  }

  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1'
  }

  const handlers = [
    http.get('/playlists', async () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => {
    server.resetHandlers()
    window.location.pathname = '/'
  })

  afterAll(() => server.close())

  it('should render minimal info', () => {
    reduxRouterRender(<AlbumCard album={album} />)

    expect(screen.getByLabelText(`default-icon-${album.title}`)).toBeInTheDocument()
    expect(screen.getByText(album.title)).toBeInTheDocument()
    expect(screen.getByText(/unknown/i)).toBeInTheDocument()
  })

  it('should render maximal info', async () => {
    const localAlbum = {
      ...album,
      imageUrl: 'something.png',
      artist: artist
    }

    reduxRouterRender(<AlbumCard album={localAlbum} />)

    expect(screen.getByRole('img', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.title })).toHaveAttribute(
      'src',
      localAlbum.imageUrl
    )
    expect(screen.getByText(localAlbum.title)).toBeInTheDocument()
    expect(screen.getByText(localAlbum.artist.name)).toBeInTheDocument()
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRouterRender(<AlbumCard album={album} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`default-icon-${album.title}`)
    })

    expect(screen.getByRole('menuitem', { name: /open drawer/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeDisabled()

    rerender(<AlbumCard album={{ ...album, artist: emptyArtist }} />)
    expect(screen.getByRole('menuitem', { name: /view artist/i })).not.toBeDisabled()
  })

  describe('on menu', () => {
    it('should open album drawer when clicking on open drawer', async () => {
      const user = userEvent.setup()

      const [_, store] = reduxRouterRender(<AlbumCard album={album} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${album.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /open drawer/i }))

      expect((store.getState() as RootState).global.albumDrawer.open).toBeTruthy()
      expect((store.getState() as RootState).global.albumDrawer.albumId).toBe(album.id)
    })

    it('should navigate to artist when clicking on view artist', async () => {
      const user = userEvent.setup()

      const localAlbum = { ...album, artist: { ...emptyArtist, id: '1' } }

      reduxRouterRender(<AlbumCard album={localAlbum} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${localAlbum.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view artist/i }))
      expect(window.location.pathname).toBe(`/artist/${localAlbum.artist.id}`)
    })

    it('should display warning modal when clicking on delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/albums/${album.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(<AlbumCard album={album} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${album.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))
    })
  })

  it('should navigate on click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<AlbumCard album={album} />)

    await user.click(screen.getByLabelText(`default-icon-${album.title}`))
    expect(window.location.pathname).toBe(`/album/${album.id}`)
  })

  it('should open artist drawer on artist name click', async () => {
    const user = userEvent.setup()

    const localAlbum = {
      ...album,
      artist: artist
    }

    const [_, store] = reduxRouterRender(<AlbumCard album={localAlbum} />)

    await user.click(screen.getByText(localAlbum.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localAlbum.artist.id)
  })
})
