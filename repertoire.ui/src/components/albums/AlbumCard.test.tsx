import Album from 'src/types/models/Album.ts'
import { reduxRouterRender, withToastify } from '../../test-utils.tsx'
import AlbumCard from './AlbumCard.tsx'
import { screen } from '@testing-library/react'
import Artist from 'src/types/models/Artist.ts'
import userEvent from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { RootState } from 'src/state/store.ts'

describe('Album Card', () => {
  const album: Album = {
    id: '1',
    title: 'Album 1',
    createdAt: '',
    updatedAt: '',
    songs: []
  }

  const artist: Artist = {
    id: '1',
    name: 'Artist 1',
    createdAt: '',
    updatedAt: '',
    albums: [],
    songs: []
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render minimal info', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRouterRender(<AlbumCard album={album} />)

    // Assert
    expect(screen.getByRole('img', { name: album.title })).toBeInTheDocument()
    expect(screen.getByText(album.title)).toBeInTheDocument()
    expect(screen.getByText(/unknown/i)).toBeInTheDocument()

    await user.pointer({ keys: '[MouseRight>]', target: screen.getByRole('img', { name: album.title }) })
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should render maximal info', async () => {
    // Arrange
    const user = userEvent.setup()

    const localAlbum = {
      ...album,
      artist: artist
    }

    // Act
    reduxRouterRender(<AlbumCard album={localAlbum} />)

    // Assert
    expect(screen.getByRole('img', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByText(localAlbum.title)).toBeInTheDocument()
    expect(screen.getByText(localAlbum.artist.name)).toBeInTheDocument()

    await user.pointer({ keys: '[MouseRight>]', target: screen.getByRole('img', { name: album.title }) })
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should display warning modal when clicking on delete', async () => {
    // Arrange
    const user = userEvent.setup()

    server.use(
      http.delete(`/albums/${album.id}`, async () => {
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    // Act
    reduxRouterRender(withToastify(<AlbumCard album={album} />))

    // Assert
    await user.pointer({ keys: '[MouseRight>]', target: screen.getByRole('img', { name: album.title }) })
    await user.click(screen.getByRole('menuitem', { name: /delete/i }))

    expect(screen.getByRole('heading', { name: /delete/i })).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(screen.getByText(`${album.title} deleted!`)).toBeInTheDocument()
  })

  it('should navigate on click', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRouterRender(<AlbumCard album={album} />)

    // Assert
    await user.click(screen.getByRole('img', { name: album.title }))
    expect(window.location.pathname).toBe(`/album/${album.id}`)
  })

  it('should open artist drawer on artist name click', async () => {
    // Arrange
    const user = userEvent.setup()

    const localAlbum = {
      ...album,
      artist: artist
    }

    // Act
    const [_, store] = reduxRouterRender(<AlbumCard album={localAlbum} />)

    // Assert
    await user.click(screen.getByText(localAlbum.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localAlbum.artist.id)
  })
})