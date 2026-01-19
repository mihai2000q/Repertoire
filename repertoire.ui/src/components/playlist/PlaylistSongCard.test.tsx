import {
  emptyAlbum,
  emptyArtist,
  emptyOrder,
  emptySong,
  reduxRouterRender,
  withToastify
} from '../../test-utils.tsx'
import PlaylistSongCard from './PlaylistSongCard.tsx'
import Song from '../../types/models/Song.ts'
import { fireEvent, screen, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Artist from '../../types/models/Artist.ts'
import Album from '../../types/models/Album.ts'
import { RootState } from '../../state/store.ts'
import { afterEach, beforeEach, expect } from 'vitest'
import { RemoveSongsFromPlaylistRequest } from '../../types/requests/PlaylistRequests.ts'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import SongProperty from '../../types/enums/properties/SongProperty.ts'
import Difficulty from '../../types/enums/Difficulty.ts'
import dayjs from 'dayjs'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../types/models/Playlist.ts'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'

describe('Playlist Song Card', () => {
  const song: Song = {
    ...emptySong,
    id: '1',
    title: 'Song 1',
    playlistTrackNo: 1,
    playlistSongId: '1-1'
  }

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
      isClickSelectionActive: false,
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

  it('should render and display minimal info', () => {
    reduxRouterRender(
      <PlaylistSongCard song={song} playlistId={''} order={emptyOrder} isDragging={false} />
    )

    expect(screen.getByText(song.playlistTrackNo)).toBeInTheDocument()
    expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
    expect(screen.getByText(song.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should render and display maximal info', () => {
    const localSong = {
      ...song,
      imageUrl: 'something.png',
      artist: artist,
      album: album
    }

    reduxRouterRender(
      <PlaylistSongCard song={localSong} playlistId={''} order={emptyOrder} isDragging={false} />
    )

    expect(screen.getByText(localSong.playlistTrackNo)).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.title })).toHaveAttribute(
      'src',
      localSong.imageUrl
    )
    expect(screen.getByText(localSong.title)).toBeInTheDocument()
    expect(screen.getByText(localSong.album.title)).toBeInTheDocument()
    expect(screen.getByText(localSong.artist.name)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it("should display song's image if the song has one, if not the album's image", () => {
    const localSong: Song = {
      ...song,
      imageUrl: 'something.png'
    }

    const [{ rerender }] = reduxRouterRender(
      <PlaylistSongCard song={localSong} playlistId={''} order={emptyOrder} isDragging={false} />
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
      <PlaylistSongCard
        song={localSongWithAlbum}
        playlistId={''}
        order={emptyOrder}
        isDragging={false}
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

      const [{ rerender }] = reduxRouterRender(
        <PlaylistSongCard song={localSong} playlistId={''} order={order} isDragging={false} />
      )

      expect(
        screen.getByText(dayjs(localSong.releaseDate).format('DD MMM YYYY'))
      ).toBeInTheDocument()

      rerender(<PlaylistSongCard song={song} playlistId={''} order={order} isDragging={false} />)

      expect(screen.getByText(/unknown/i)).toBeInTheDocument()
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
        <PlaylistSongCard song={localSong} playlistId={''} order={order} isDragging={false} />
      )

      expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()

      rerender(<PlaylistSongCard song={song} playlistId={''} order={order} isDragging={false} />)

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
        <PlaylistSongCard song={localSong} playlistId={''} order={order} isDragging={false} />
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
        <PlaylistSongCard song={localSong} playlistId={''} order={order} isDragging={false} />
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
        <PlaylistSongCard song={localSong} playlistId={''} order={order} isDragging={false} />
      )

      expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    })

    it('should display the last played date, when it is last  played', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.LastPlayed
      }

      const localSong = {
        ...song,
        lastTimePlayed: '2024-10-12T10:30'
      }

      const [{ rerender }] = reduxRouterRender(
        <PlaylistSongCard song={localSong} playlistId={''} order={order} isDragging={false} />
      )

      expect(
        screen.getByText(dayjs(localSong.lastTimePlayed).format('DD MMM YYYY'))
      ).toBeInTheDocument()

      rerender(<PlaylistSongCard song={song} playlistId={''} order={order} isDragging={false} />)

      expect(screen.getByText(/never/i)).toBeInTheDocument()
    })
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRouterRender(
      <PlaylistSongCard song={song} playlistId={''} order={emptyOrder} isDragging={false} />
    )
    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-card-${song.title}`)
    })
    expectDefaultMenu()

    rerender(
      <PlaylistSongCard
        song={{
          ...song,
          artist: emptyArtist,
          album: emptyAlbum,
          songsterrLink: 'some-link',
          youtubeLink: 'https://www.youtube.com/watch?v=tAGnKpE4NCI'
        }}
        playlistId={''}
        order={emptyOrder}
        isDragging={false}
      />
    )
    await expectFullMenu()
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRouterRender(
      <PlaylistSongCard song={song} playlistId={''} order={emptyOrder} isDragging={false} />
    )
    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expectDefaultMenu()

    rerender(
      <PlaylistSongCard
        song={{
          ...song,
          artist: emptyArtist,
          album: emptyAlbum,
          songsterrLink: 'some-link',
          youtubeLink: 'https://www.youtube.com/watch?v=tAGnKpE4NCI'
        }}
        playlistId={''}
        order={emptyOrder}
        isDragging={false}
      />
    )
    await expectFullMenu()
  })

  function expectDefaultMenu() {
    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /partial rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove from playlist/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeDisabled()
  }

  async function expectFullMenu() {
    const user = userEvent.setup()

    expect(screen.getByRole('menuitem', { name: /view artist/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /open links/i })).not.toBeDisabled()
    await user.hover(screen.getByRole('menuitem', { name: /open links/i }))
    expect(screen.getByRole('menuitem', { name: /songsterr/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /youtube/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /songsterr/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /youtube/i })).not.toBeDisabled()
  }

  describe('on menu', () => {
    it('should navigate to song when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <PlaylistSongCard song={song} playlistId={''} order={emptyOrder} isDragging={false} />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))

      expect(window.location.pathname).toBe(`/song/${song.id}`)
    })

    it('should navigate to artist when clicking on view artist', async () => {
      const user = userEvent.setup()

      const localSong = { ...song, artist: emptyArtist }

      reduxRouterRender(
        <PlaylistSongCard song={localSong} playlistId={''} order={emptyOrder} isDragging={false} />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view artist/i }))

      expect(window.location.pathname).toBe(`/artist/${localSong.artist.id}`)
    })

    it('should navigate to album when clicking on view album', async () => {
      const user = userEvent.setup()

      const localSong = { ...song, album: emptyAlbum }

      reduxRouterRender(
        <PlaylistSongCard song={localSong} playlistId={''} order={emptyOrder} isDragging={false} />
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
        <PlaylistSongCard song={localSong} playlistId={''} order={emptyOrder} isDragging={false} />
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
        <PlaylistSongCard song={localSong} playlistId={''} order={emptyOrder} isDragging={false} />
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

    it('should display warning modal and remove, when clicking on remove', async () => {
      const user = userEvent.setup()

      let capturedRequest: RemoveSongsFromPlaylistRequest
      server.use(
        http.put('/playlists/songs/remove', async (req) => {
          capturedRequest = (await req.request.json()) as RemoveSongsFromPlaylistRequest
          return HttpResponse.json()
        })
      )

      const playlistId = 'some-id'

      reduxRouterRender(
        withToastify(
          <PlaylistSongCard
            song={song}
            playlistId={playlistId}
            order={emptyOrder}
            isDragging={false}
          />
        )
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /remove/i }))

      expect(await screen.findByRole('dialog', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(capturedRequest.id).toBe(playlistId)
      expect(capturedRequest.playlistSongIds).toStrictEqual([song.playlistSongId])
      expect(await screen.findByText(new RegExp(`${song.title} removed`, 'i'))).toBeInTheDocument()
    })
  })

  it('should open album drawer on album title click', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      album: album
    }

    const [_, store] = reduxRouterRender(
      <PlaylistSongCard song={localSong} playlistId={''} order={emptyOrder} isDragging={false} />
    )

    await user.click(screen.getByText(localSong.album.title))

    expect((store.getState() as RootState).global.albumDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.albumDrawer.albumId).toBe(localSong.album.id)
  })

  it('should open artist drawer on artist name click', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      artist: artist
    }

    const [_, store] = reduxRouterRender(
      <PlaylistSongCard song={localSong} playlistId={''} order={emptyOrder} isDragging={false} />
    )

    await user.click(screen.getByText(localSong.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localSong.artist.id)
  })

  it('should disable context menu and more menu when click selection is active', async () => {
    const user = userEvent.setup()

    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [],
      isClickSelectionActive: true,
      clearSelection: vi.fn()
    })

    reduxRouterRender(
      <PlaylistSongCard song={song} playlistId={''} order={emptyOrder} isDragging={false} />
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
        <PlaylistSongCard song={song} playlistId={''} order={emptyOrder} isDragging={false} />
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
        <PlaylistSongCard song={song} playlistId={''} order={emptyOrder} isDragging={false} />
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
        <PlaylistSongCard song={song} playlistId={''} order={emptyOrder} isDragging={false} />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))

      expect(screen.getByLabelText(`song-card-${song.title}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when is dragging', () => {
      reduxRouterRender(
        <PlaylistSongCard song={song} playlistId={''} order={emptyOrder} isDragging={true} />
      )

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
        selectedIds: [song.playlistSongId],
        isClickSelectionActive: true,
        clearSelection: vi.fn()
      })

      reduxRouterRender(
        <PlaylistSongCard song={song} playlistId={''} order={emptyOrder} isDragging={false} />
      )

      expect(screen.getByLabelText(`song-card-${song.title}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })
  })
})
