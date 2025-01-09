import Album from 'src/types/models/Album.ts'
import { reduxRouterRender, withToastify } from '../../test-utils.tsx'
import AlbumHeaderCard from './AlbumHeaderCard.tsx'
import userEvent from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import Song from 'src/types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import Artist from 'src/types/models/Artist.ts'
import { RootState } from 'src/state/store.ts'

describe('Album Header Card', () => {
  const emptySong: Song = {
    id: '',
    title: '',
    description: '',
    isRecorded: false,
    rehearsals: 0,
    confidence: 0,
    progress: 0,
    sections: [],
    createdAt: '',
    updatedAt: ''
  }

  const album: Album = {
    id: '1',
    title: 'Album 1',
    createdAt: '',
    updatedAt: '',
    releaseDate: null,
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


  it('should render and display minimal info when the album is not unknown', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRouterRender(<AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />)

    // Assert
    expect(screen.getByRole('img', { name: album.title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: album.title })).toBeInTheDocument()
    expect(screen.getByText('0 songs')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /info/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should render and display maximal info when the album is not unknown', async () => {
    // Arrange
    const user = userEvent.setup()

    const localAlbum: Album = {
      ...album,
      releaseDate: '2024-10-21T10:30:00',
      artist: artist,
      songs: [
        {
          ...emptySong,
          id: '1',
          title: 'Song 1'
        },
        {
          ...emptySong,
          id: '2',
          title: 'Song 2'
        }
      ]
    }

    // Act
    reduxRouterRender(<AlbumHeaderCard album={localAlbum} isUnknownAlbum={false} songsTotalCount={undefined} />)

    // Assert
    expect(screen.getByRole('img', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.artist.name })).toBeInTheDocument()
    expect(screen.getByText(localAlbum.artist.name)).toBeInTheDocument()
    expect(screen.getByText(/2024/)).toBeInTheDocument()
    expect(screen.getByText(`${localAlbum.songs.length} songs`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()

    await user.hover(screen.getByText(/2024/))
    expect(await screen.findByText(/21 October 2024/)).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /info/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should render and display info when the album is unknown', async () => {
    // Arrange
    const songsTotalCount = 10

    // Act
    reduxRouterRender(<AlbumHeaderCard album={undefined} isUnknownAlbum={true} songsTotalCount={songsTotalCount} />)

    // Assert
    expect(screen.getByRole('img', { name: 'unknown-album' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /unknown/i })).toBeInTheDocument()
    expect(screen.getByText(`${songsTotalCount} songs`)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'more-menu' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'edit-header' })).not.toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display info modal', async () => {
      // Arrange
      const user = userEvent.setup()

      // Act
      reduxRouterRender(<AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />)

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /info/i }))

      expect(screen.getByRole('heading', { name: /album info/i })).toBeInTheDocument()
    })

    it('should display edit header modal', async () => {
      // Arrange
      const user = userEvent.setup()

      // Act
      reduxRouterRender(<AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />)

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(screen.getByRole('heading', { name: /edit album header/i })).toBeInTheDocument()
    })

    it('should display warning modal and delete album', async () => {
      // Arrange
      const user = userEvent.setup()

      server.use(
        http.delete(`/albums/${album.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      // Act
      reduxRouterRender(
        withToastify(<AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />)
      )

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(screen.getByRole('heading', { name: /delete album/i })).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(window.location.pathname).toBe('/albums')
      expect(screen.getByText(`${album.title} deleted!`)).toBeInTheDocument()
    })
  })

  it('should display edit header modal from edit button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRouterRender(<AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />)

    // Assert
    await user.click(screen.getByRole('button', { name: 'edit-header' }))

    expect(screen.getByRole('heading', { name: /edit album header/i })).toBeInTheDocument()
  })

  it('should open artist drawer on artist click', async () => {
    // Arrange
    const user = userEvent.setup()

    const localAlbum: Album = {
      ...album,
      artist: artist
    }

    // Act
    const [_, store] = reduxRouterRender(<AlbumHeaderCard album={localAlbum} isUnknownAlbum={false} songsTotalCount={undefined} />)

    // Assert
    await user.click(screen.getByText(localAlbum.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localAlbum.artist.id)
  })
})
