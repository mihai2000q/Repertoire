import {
  emptyAlbum,
  emptyOrder,
  emptySong,
  reduxRouterRender,
  withToastify
} from '../../test-utils.tsx'
import ArtistSongCard from './ArtistSongCard.tsx'
import Song from '../../types/models/Song.ts'
import { fireEvent, screen, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Album from 'src/types/models/Album.ts'
import { RootState } from '../../state/store.ts'
import dayjs from 'dayjs'
import SongProperty from '../../types/enums/properties/SongProperty.ts'
import Difficulty from '../../types/enums/Difficulty.ts'
import { beforeEach, expect } from 'vitest'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { RemoveSongsFromArtistRequest } from '../../types/requests/ArtistRequests.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../types/models/Playlist.ts'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'

describe('Artist Song Card', () => {
  const song: Song = {
    ...emptySong,
    id: '1',
    title: 'Song 1'
  }

  const album: Album = {
    ...emptyAlbum,
    id: '1',
    title: 'Album 1'
  }

  const handlers = [
    http.get('/playlists', async () => {
      const response: WithTotalCountResponse<Playlist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [],
      clearSelection: vi.fn()
    })
  })

  afterEach(() => {
    vi.restoreAllMocks()
    server.resetHandlers()
    window.location.pathname = '/'
  })

  beforeAll(() => {
    vi.mock('../../context/ClickSelectContext', () => ({
      useClickSelect: vi.fn()
    }))
    server.listen()
  })

  afterAll(() => server.close())

  it('should render and display minimal information', async () => {
    reduxRouterRender(
      <ArtistSongCard song={song} artistId={''} isUnknownArtist={false} order={emptyOrder} />
    )

    expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
    expect(screen.getByText(song.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should render and display maximal information', async () => {
    const localSong: Song = {
      ...song,
      imageUrl: 'something.png',
      album: album
    }

    reduxRouterRender(
      <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={emptyOrder} />
    )

    expect(screen.getByRole('img', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.title })).toHaveAttribute(
      'src',
      localSong.imageUrl
    )
    expect(screen.getByText(localSong.title)).toBeInTheDocument()
    expect(screen.getByText(localSong.album.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it("should display song's image if the song has one, if not the album's image", () => {
    const localSong: Song = {
      ...song,
      imageUrl: 'something.png'
    }

    const [{ rerender }] = reduxRouterRender(
      <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={emptyOrder} />
    )

    expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', localSong.imageUrl)

    const localSongWithAlbum: Song = {
      ...song,
      album: {
        ...emptyAlbum,
        imageUrl: 'something-album.png'
      }
    }

    rerender(
      <ArtistSongCard
        song={localSongWithAlbum}
        artistId={''}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    expect(screen.getByRole('img', { name: song.title })).toHaveAttribute(
      'src',
      localSongWithAlbum.album.imageUrl
    )
  })

  describe('on order property change', () => {
    it('should display the release date, when it is release date', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.ReleaseDate
      }

      const localSong = {
        ...song,
        releaseDate: '2024-10-12T10:30'
      }

      reduxRouterRender(
        <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={order} />
      )

      expect(
        screen.getByText(dayjs(localSong.releaseDate).format('DD MMM YYYY'))
      ).toBeInTheDocument()
    })

    it('should display the difficulty bar, when it is difficulty', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.Difficulty
      }

      const localSong = {
        ...song,
        difficulty: Difficulty.Easy
      }

      const [{ rerender }] = reduxRouterRender(
        <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={order} />
      )

      expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()

      rerender(<ArtistSongCard song={song} artistId={''} isUnknownArtist={false} order={order} />)

      expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()
    })

    it('should display the rehearsals, when it is rehearsals', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.Rehearsals
      }

      const localSong = {
        ...song,
        rehearsals: 34
      }

      reduxRouterRender(
        <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={order} />
      )

      expect(screen.getAllByText(localSong.rehearsals)).toHaveLength(2) // the one visible and the one in the tooltip
    })

    it('should display the confidence bar, when it is confidence', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.Confidence
      }

      const localSong = {
        ...song,
        confidence: 23
      }

      reduxRouterRender(
        <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={order} />
      )

      expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    })

    it('should display the progress bar, when it is progress', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.Progress
      }

      const localSong = {
        ...song,
        progress: 123
      }

      reduxRouterRender(
        <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={order} />
      )

      expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    })

    it('should display the last played date, when it is last played', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.LastPlayed
      }

      const localSong = {
        ...song,
        lastTimePlayed: '2024-10-12T10:30'
      }

      const [{ rerender }] = reduxRouterRender(
        <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={order} />
      )

      expect(
        screen.getByText(dayjs(localSong.lastTimePlayed).format('DD MMM YYYY'))
      ).toBeInTheDocument()

      rerender(<ArtistSongCard song={song} artistId={''} isUnknownArtist={false} order={order} />)

      expect(screen.getByText(/never/i)).toBeInTheDocument()
    })
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRouterRender(
      <ArtistSongCard song={song} artistId={''} isUnknownArtist={false} order={emptyOrder} />
    )
    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-card-${song.title}`)
    })
    expectDefaultMenu()

    rerender(
      <ArtistSongCard
        song={{
          ...song,
          album: emptyAlbum,
          youtubeLink: 'https://www.youtube.com/watch?v=tAGnKpE4NCI',
          songsterrLink: 'some-link'
        }}
        artistId={''}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )
    await expectFullMenu()
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRouterRender(
      <ArtistSongCard song={song} artistId={''} isUnknownArtist={false} order={emptyOrder} />
    )
    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expectDefaultMenu()

    rerender(
      <ArtistSongCard
        song={{
          ...song,
          album: emptyAlbum,
          youtubeLink: 'https://www.youtube.com/watch?v=tAGnKpE4NCI',
          songsterrLink: 'some-link'
        }}
        artistId={''}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )
    await expectFullMenu()
  })

  function expectDefaultMenu() {
    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /partial rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove from artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeDisabled()
  }

  async function expectFullMenu() {
    const user = userEvent.setup()

    expect(screen.getByRole('menuitem', { name: /view album/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /open links/i })).not.toBeDisabled()
    await user.hover(screen.getByRole('menuitem', { name: /open links/i }))
    expect(screen.getByRole('menuitem', { name: /songsterr/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /youtube/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /songsterr/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /youtube/i })).not.toBeDisabled()
  }

  it('should display less information on the menu when the artist is unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistSongCard song={song} artistId={''} isUnknownArtist={true} order={emptyOrder} />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.queryByRole('menuitem', { name: /remove from artist/i })).not.toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to song when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistSongCard song={song} artistId={''} isUnknownArtist={false} order={emptyOrder} />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))

      expect(window.location.pathname).toBe(`/song/${song.id}`)
    })

    it('should navigate to album when clicking on view album', async () => {
      const user = userEvent.setup()

      const localSong = { ...song, album: emptyAlbum }

      reduxRouterRender(
        <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={emptyOrder} />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view album/i }))

      expect(window.location.pathname).toBe(`/album/${localSong.album.id}`)
    })

    it('should open youtube when clicking on open links => open youtube', async () => {
      const user = userEvent.setup()

      const localSong = { ...song, youtubeLink: 'https://www.youtube.com/watch?v=tAGnKpE4NCI' }

      server.use(
        http.get(
          localSong.youtubeLink
            .replace('youtube', 'youtube-nocookie')
            .replace('watch?v=', 'embed/'),
          () => {
            return HttpResponse.json({ message: 'it worked' })
          }
        )
      )

      reduxRouterRender(
        <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={emptyOrder} />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.hover(screen.getByRole('menuitem', { name: /open links/i }))
      fireEvent.click(screen.getByRole('menuitem', { name: /youtube/i }))
      expect(await screen.findByRole('dialog', { name: song.title })).toBeInTheDocument()
    })

    it('should be able to open songsterr in browser when clicking on open links => open songsterr', async () => {
      const user = userEvent.setup()

      const localSong = { ...song, songsterrLink: 'some-link' }

      reduxRouterRender(
        <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={emptyOrder} />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.hover(screen.getByRole('menuitem', { name: /open links/i }))
      expect(
        within(screen.getByRole('link', { name: /songsterr/i })).getByRole('menuitem', {
          name: /songsterr/i
        })
      ).toBeInTheDocument()
      expect(screen.getByRole('link', { name: /songsterr/i })).toBeExternalLink(
        localSong.songsterrLink
      )
    })

    it('should display warning modal and remove song from artist, when clicking on remove from artist', async () => {
      const user = userEvent.setup()

      let capturedRequest: RemoveSongsFromArtistRequest
      server.use(
        http.put('/artists/remove-songs', async (req) => {
          capturedRequest = (await req.request.json()) as RemoveSongsFromArtistRequest
          return HttpResponse.json()
        })
      )

      const artistId = 'some-artist-id'

      reduxRouterRender(
        withToastify(
          <ArtistSongCard
            song={song}
            artistId={artistId}
            isUnknownArtist={false}
            order={emptyOrder}
          />
        )
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /remove from artist/i }))

      expect(await screen.findByRole('dialog', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(capturedRequest).toStrictEqual({
        id: artistId,
        songIds: [song.id]
      })
      expect(await screen.findByText(new RegExp(`${song.title} removed`, 'i'))).toBeInTheDocument()
    })

    it('should display warning modal and delete song, when clicking on delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/songs/${song.id}`, () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(
        <ArtistSongCard song={song} artistId={''} isUnknownArtist={false} order={emptyOrder} />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete song/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete song/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))
    })
  })

  it('should open album drawer on album title click', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      album: album
    }

    const [_, store] = reduxRouterRender(
      <ArtistSongCard song={localSong} artistId={''} isUnknownArtist={false} order={emptyOrder} />
    )

    await user.click(screen.getByText(localSong.album.title))

    expect((store.getState() as RootState).global.albumDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.albumDrawer.albumId).toBe(localSong.album.id)
  })

  it('should disable context menu and more menu when click selection is active', async () => {
    const user = userEvent.setup()

    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: ['anything'],
      clearSelection: vi.fn()
    })

    reduxRouterRender(
      <ArtistSongCard artistId={''} song={song} isUnknownArtist={false} order={emptyOrder} />
    )

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`default-icon-${song.title}`)
    })
    expect(screen.queryByRole('menu')).not.toBeInTheDocument()

    expect(screen.getByRole('button', { name: 'more-menu' })).toBeDisabled()
  })

  describe('should be selected', () => {
    it('when avatar is hovered', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistSongCard artistId={''} song={song} isUnknownArtist={false} order={emptyOrder} />
      )

      await user.hover(screen.getByLabelText(`default-icon-${song.title}`))

      expect(screen.getByLabelText(`song-card-${song.title}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when context menu is open', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistSongCard artistId={''} song={song} isUnknownArtist={false} order={emptyOrder} />
      )

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${song.title}`)
      })

      expect(screen.getByLabelText(`song-card-${song.title}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when more menu is open', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistSongCard artistId={''} song={song} isUnknownArtist={false} order={emptyOrder} />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))

      expect(screen.getByLabelText(`song-card-${song.title}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when click selected (part of the selected ids)', () => {
      vi.mocked(useClickSelect).mockReturnValue({
        selectables: [],
        addSelectable: vi.fn(),
        removeSelectable: vi.fn(),
        selectedIds: [song.id],
        clearSelection: vi.fn()
      })

      reduxRouterRender(
        <ArtistSongCard artistId={''} song={song} isUnknownArtist={false} order={emptyOrder} />
      )

      expect(screen.getByLabelText(`song-card-${song.title}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })
  })
})
