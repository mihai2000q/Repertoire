import {
  emptyAlbum,
  emptyArtist,
  emptyOrder,
  emptySong,
  reduxRouterRender
} from '../../test-utils.tsx'
import PlaylistSongCard from './PlaylistSongCard.tsx'
import Song from '../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Artist from '../../types/models/Artist.ts'
import Album from '../../types/models/Album.ts'
import { RootState } from '../../state/store.ts'
import { afterEach, expect } from 'vitest'
import { RemoveSongsFromPlaylistRequest } from '../../types/requests/PlaylistRequests.ts'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import SongProperty from '../../types/enums/SongProperty.ts'
import Difficulty from '../../types/enums/Difficulty.ts'
import dayjs from 'dayjs'

describe('Playlist Song Card', () => {
  const song: Song = {
    ...emptySong,
    id: '1',
    title: 'Song 1',
    playlistTrackNo: 1
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

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => {
    server.resetHandlers()
    window.location.pathname = '/'
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
      expect(screen.getAllByText(localSong.rehearsals)[0]).toBeVisible()
      expect(screen.getAllByText(localSong.rehearsals)[1]).not.toBeVisible()
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

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /partial rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeDisabled()

    rerender(
      <PlaylistSongCard
        song={{ ...song, artist: emptyArtist, album: emptyAlbum }}
        playlistId={''}
        order={emptyOrder}
        isDragging={false}
      />
    )
    expect(screen.getByRole('menuitem', { name: /view artist/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).not.toBeDisabled()
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRouterRender(
      <PlaylistSongCard song={song} playlistId={''} order={emptyOrder} isDragging={false} />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /partial rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeDisabled()

    rerender(
      <PlaylistSongCard
        song={{ ...song, artist: emptyArtist, album: emptyAlbum }}
        playlistId={''}
        order={emptyOrder}
        isDragging={false}
      />
    )
    expect(screen.getByRole('menuitem', { name: /view artist/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).not.toBeDisabled()
  })

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

    it('should display warning modal and remove, when clicking on remove', async () => {
      const user = userEvent.setup()

      let capturedRequest: RemoveSongsFromPlaylistRequest
      server.use(
        http.put('/playlists/remove-songs', async (req) => {
          capturedRequest = (await req.request.json()) as RemoveSongsFromPlaylistRequest
          return HttpResponse.json()
        })
      )

      const playlistId = 'some-id'

      reduxRouterRender(
        <PlaylistSongCard
          song={song}
          playlistId={playlistId}
          order={emptyOrder}
          isDragging={false}
        />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /remove/i }))

      expect(await screen.findByRole('dialog', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(capturedRequest.id).toBe(playlistId)
      expect(capturedRequest.songIds).toStrictEqual([song.id])
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
})
