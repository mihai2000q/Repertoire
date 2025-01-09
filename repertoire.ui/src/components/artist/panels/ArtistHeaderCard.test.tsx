import { reduxRouterRender, withToastify } from '../../../test-utils.tsx'
import ArtistHeaderCard from './ArtistHeaderCard.tsx'
import userEvent from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import Song from 'src/types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import Artist from 'src/types/models/Artist.ts'
import { RootState } from 'src/state/store.ts'

describe('Artist Header Card', () => {
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

  const artist: Artist = {
    id: '1',
    name: 'Artist 1',
    createdAt: '',
    updatedAt: '',
    albums: [],
    songs: []
  }

  const albumsTotalCount = 10
  const songsTotalCount = 20

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())


  it('should render and display minimal info when the artist is not unknown', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRouterRender(
      <ArtistHeaderCard
        artist={artist}
        albumsTotalCount={albumsTotalCount}
        songsTotalCount={songsTotalCount}
        isUnknownArtist={false}
      />
    )

    // Assert
    expect(screen.getByRole('img', { name: artist.name })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: artist.name })).toBeInTheDocument()
    expect(screen.getByText(`${albumsTotalCount} albums • ${songsTotalCount} songs`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /info/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should render and display maximal info when the artist is not unknown', async () => {
    // Arrange
    const user = userEvent.setup()

    const localArtist: Artist = {
      ...artist,
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

    const albumsTotalCount = 1
    const songsTotalCount = 1

    // Act
    reduxRouterRender(
      <ArtistHeaderCard
        artist={localArtist}
        albumsTotalCount={albumsTotalCount}
        songsTotalCount={songsTotalCount}
        isUnknownArtist={false}
      />
    )

    // Assert
    expect(screen.getByRole('img', { name: localArtist.name })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: localArtist.name })).toBeInTheDocument()
    expect(screen.getByText(`${albumsTotalCount} album • ${songsTotalCount} song`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /info/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should render and display info when the artist is unknown', async () => {
    reduxRouterRender(
      <ArtistHeaderCard
        artist={undefined}
        albumsTotalCount={albumsTotalCount}
        songsTotalCount={songsTotalCount}
        isUnknownArtist={true}
      />
    )

    expect(screen.getByRole('img', { name: 'unknown-artist' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /unknown/i })).toBeInTheDocument()
    expect(screen.getByText(`${albumsTotalCount} albums • ${songsTotalCount} songs`)).toBeInTheDocument()

    expect(screen.queryByRole('button', { name: 'more-menu' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'edit-header' })).not.toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display info modal', async () => {
      // Arrange
      const user = userEvent.setup()

      // Act
      reduxRouterRender(
        <ArtistHeaderCard
          artist={artist}
          albumsTotalCount={undefined}
          songsTotalCount={undefined}
          isUnknownArtist={false}
        />
      )

      // Assert
      await user.hover(screen.getByLabelText('header-panel-card'))
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /info/i }))

      expect(screen.getByRole('heading', { name: /artist info/i })).toBeInTheDocument()
    })

    it('should display edit header modal', async () => {
      // Arrange
      const user = userEvent.setup()

      // Act
      reduxRouterRender(
        <ArtistHeaderCard
          artist={artist}
          albumsTotalCount={undefined}
          songsTotalCount={undefined}
          isUnknownArtist={false}
        />
      )

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(screen.getByRole('heading', { name: /edit artist header/i })).toBeInTheDocument()
    })

    it('should display warning modal and delete artist', async () => {
      // Arrange
      const user = userEvent.setup()

      server.use(
        http.delete(`/artists/${artist.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      // Act
      reduxRouterRender(
        withToastify(
          <ArtistHeaderCard
            artist={artist}
            albumsTotalCount={undefined}
            songsTotalCount={undefined}
            isUnknownArtist={false}
          />
        )
      )

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(screen.getByRole('heading', { name: /delete artist/i })).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(window.location.pathname).toBe('/artists')
      expect(screen.getByText(`${artist.name} deleted!`)).toBeInTheDocument()
    })
  })

  it('should display edit header modal from edit button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRouterRender(
      <ArtistHeaderCard
        artist={artist}
        albumsTotalCount={undefined}
        songsTotalCount={undefined}
        isUnknownArtist={false}
      />
    )

    // Assert
    await user.click(screen.getByRole('button', { name: 'edit-header' }))

    expect(screen.getByRole('heading', { name: /edit artist header/i })).toBeInTheDocument()
  })
})
