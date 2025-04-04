import Album from 'src/types/models/Album.ts'
import {emptyAlbum, emptyArtist, reduxRouterRender, withToastify} from '../../test-utils.tsx'
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
    title: 'Album 1',
  }

  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1',
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

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
    expect(screen.getByRole('img', { name: localAlbum.title })).toHaveAttribute('src', localAlbum.imageUrl)
    expect(screen.getByText(localAlbum.title)).toBeInTheDocument()
    expect(screen.getByText(localAlbum.artist.name)).toBeInTheDocument()
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<AlbumCard album={album} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`default-icon-${album.title}`)
    })
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display warning modal when clicking on delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/albums/${album.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(withToastify(<AlbumCard album={album} />))

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${album.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(screen.getByText(`${album.title} deleted!`)).toBeInTheDocument()
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
