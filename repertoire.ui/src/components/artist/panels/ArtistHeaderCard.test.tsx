import { emptyArtist, emptySong, reduxRouterRender } from '../../../test-utils.tsx'
import ArtistHeaderCard from './ArtistHeaderCard.tsx'
import userEvent from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import Artist from 'src/types/models/Artist.ts'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../../types/models/Playlist.ts'

describe('Artist Header Card', () => {
  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1',
    bandMembers: [
      {
        id: '1',
        name: 'Member 1',
        roles: []
      },
      {
        id: '2',
        name: 'Member 2',
        roles: []
      }
    ]
  }

  const albumsTotalCount = 10
  const songsTotalCount = 20

  const handlers = [
    http.get('/playlists', async () => {
      const response: WithTotalCountResponse<Playlist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display minimal info when the artist is not unknown', () => {
    reduxRouterRender(
      <ArtistHeaderCard
        artist={artist}
        albumsTotalCount={albumsTotalCount}
        songsTotalCount={songsTotalCount}
        isUnknownArtist={false}
      />
    )

    expect(screen.getByLabelText(`default-icon-${artist.name}`)).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: artist.name })).toBeInTheDocument()
    expect(
      screen.getByText(`${albumsTotalCount} albums • ${songsTotalCount} songs`)
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()
  })

  it('should render and display maximal info when the artist is not unknown', () => {
    const localArtist: Artist = {
      ...artist,
      isBand: true,
      imageUrl: 'something.png',
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

    reduxRouterRender(
      <ArtistHeaderCard
        artist={localArtist}
        albumsTotalCount={albumsTotalCount}
        songsTotalCount={songsTotalCount}
        isUnknownArtist={false}
      />
    )

    expect(screen.getByRole('img', { name: localArtist.name })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: artist.name })).toHaveAttribute('src', artist.imageUrl)
    expect(screen.getByRole('heading', { name: localArtist.name })).toBeInTheDocument()
    expect(
      screen.getByText(
        `${artist.bandMembers.length} members • ${albumsTotalCount} album • ${songsTotalCount} song`
      )
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()
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

    expect(screen.getByLabelText('icon-unknown-artist')).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /unknown/i })).toBeInTheDocument()
    expect(
      screen.getByText(`${albumsTotalCount} albums • ${songsTotalCount} songs`)
    ).toBeInTheDocument()

    expect(screen.queryByRole('button', { name: 'more-menu' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'edit-header' })).not.toBeInTheDocument()
  })

  it('should display image modal, when clicking the image', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistHeaderCard
        artist={{ ...artist, imageUrl: 'something.png' }}
        isUnknownArtist={false}
        songsTotalCount={undefined}
        albumsTotalCount={undefined}
      />
    )

    await user.click(screen.getByRole('img', { name: artist.name }))
    expect(await screen.findByRole('dialog', { name: artist.name + '-image' })).toBeInTheDocument()
  })

  it('should display menu on click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistHeaderCard
        artist={artist}
        albumsTotalCount={albumsTotalCount}
        songsTotalCount={songsTotalCount}
        isUnknownArtist={false}
      />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /info/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display info modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistHeaderCard
          artist={artist}
          albumsTotalCount={undefined}
          songsTotalCount={undefined}
          isUnknownArtist={false}
        />
      )

      await user.hover(screen.getByLabelText('header-panel-card'))
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /info/i }))

      expect(await screen.findByRole('dialog', { name: /artist info/i })).toBeInTheDocument()
    })

    it('should display edit header modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistHeaderCard
          artist={artist}
          albumsTotalCount={undefined}
          songsTotalCount={undefined}
          isUnknownArtist={false}
        />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(await screen.findByRole('dialog', { name: /edit artist header/i })).toBeInTheDocument()
    })

    it('should display warning modal and delete artist', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/artists/${artist.id}`, () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(
        <ArtistHeaderCard
          artist={artist}
          albumsTotalCount={undefined}
          songsTotalCount={undefined}
          isUnknownArtist={false}
        />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete artist/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))
      expect(window.location.pathname).toBe('/artists')
    })
  })

  it('should display edit header modal from edit button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistHeaderCard
        artist={artist}
        albumsTotalCount={undefined}
        songsTotalCount={undefined}
        isUnknownArtist={false}
      />
    )

    await user.click(screen.getByRole('button', { name: 'edit-header' }))

    expect(await screen.findByRole('dialog', { name: /edit artist header/i })).toBeInTheDocument()
  })
})
