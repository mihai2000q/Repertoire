import { emptyAlbum, emptyArtist, emptySong, reduxRouterRender } from '../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { afterEach, expect } from 'vitest'
import { RootState } from '../../../state/store.ts'
import Song from '../../../types/models/Song.ts'
import HomeRecentlyPlayedSongCard from './HomeRecentlyPlayedSongCard.tsx'
import dayjs from 'dayjs'
import { userEvent } from '@testing-library/user-event'

describe('Home Recently Played Song Card', () => {
  const song: Song = {
    ...emptySong,
    id: '1',
    imageUrl: 'something.png',
    title: 'Song 1',
    lastTimePlayed: '2024-12-14T10:30',
    progress: 500,
    album: {
      ...emptyAlbum,
      id: '1',
      imageUrl: 'something.png'
    },
    artist: {
      ...emptyArtist,
      id: '1',
      name: 'another-artist'
    },
    songsterrLink: 'some-link'
  }

  afterEach(() => (window.location.pathname = '/'))

  it('should render', () => {
    const localSong = {
      ...song,
      imageUrl: null,
      album: null,
      artist: null
    }

    reduxRouterRender(<HomeRecentlyPlayedSongCard song={localSong} />)

    expect(screen.getByText(localSong.title)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByLabelText(`default-icon-${localSong.title}`)).toBeInTheDocument()
    expect(screen.getByText(dayjs(localSong.lastTimePlayed).format('DD MMM'))).toBeInTheDocument()
  })

  it('should render with artist', () => {
    reduxRouterRender(<HomeRecentlyPlayedSongCard song={song} />)

    expect(screen.getByText(song.artist.name)).toBeInTheDocument()
  })

  it("should display song's image if the song has one, if not the album's image", () => {
    const [{ rerender }] = reduxRouterRender(<HomeRecentlyPlayedSongCard song={song} />)

    expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', song.imageUrl)

    const localSongWithAlbum: Song = {
      ...song,
      imageUrl: null,
      album: {
        ...emptyAlbum,
        imageUrl: 'something-album.png'
      }
    }

    rerender(<HomeRecentlyPlayedSongCard song={localSongWithAlbum} />)

    expect(screen.getByRole('img', { name: localSongWithAlbum.title })).toHaveAttribute(
      'src',
      localSongWithAlbum.album.imageUrl
    )
  })

  it('should open song drawer on clicking', async () => {
    const user = userEvent.setup()

    const [_, store] = reduxRouterRender(<HomeRecentlyPlayedSongCard song={song} />)

    await user.click(await screen.findByText(song.title))

    expect((store.getState() as RootState).global.songDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.songDrawer.songId).toBe(song.id)
  })

  it("should open artist drawer by clicking the artist", async () => {
    const user = userEvent.setup()

    const [_, store] = reduxRouterRender(<HomeRecentlyPlayedSongCard song={song} />)

    await user.click(await screen.findByText(song.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(song.artist.id)
  })

  it('should display menu on song right click', async () => {
    const user = userEvent.setup()

    const localSong: Song = {
      ...song,
      artist: null,
      album: null,
      songsterrLink: null
    }

    const [{ rerender }] = reduxRouterRender(<HomeRecentlyPlayedSongCard song={localSong} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: await screen.findByRole('img', { name: localSong.title })
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeDisabled()

    rerender(<HomeRecentlyPlayedSongCard song={song} />)
    await user.pointer({
      keys: '[MouseRight>]',
      target: await screen.findByRole('img', { name: song.title })
    })
    expect(screen.getByRole('menuitem', { name: /view artist/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /open links/i })).not.toBeDisabled()
  })

  describe('on menu', () => {
    it('should navigate to song when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<HomeRecentlyPlayedSongCard song={song} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: await screen.findByRole('img', { name: song.title })
      })
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/song/${song.id}`)
    })

    it('should navigate to artist when clicking on view artist', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<HomeRecentlyPlayedSongCard song={song} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: await screen.findByRole('img', { name: song.title })
      })
      await user.click(screen.getByRole('menuitem', { name: /view artist/i }))
      expect(window.location.pathname).toBe(`/artist/${song.artist.id}`)
    })

    it('should navigate to album when clicking on view album', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<HomeRecentlyPlayedSongCard song={song} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: await screen.findByRole('img', { name: song.title })
      })
      await user.click(screen.getByRole('menuitem', { name: /view album/i }))
      expect(window.location.pathname).toBe(`/album/${song.album.id}`)
    })
  })
})
