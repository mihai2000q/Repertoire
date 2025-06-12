import {
  emptyAlbum,
  emptyArtist,
  emptySong,
  reduxRouterRender,
  withToastify
} from '../../../test-utils.tsx'
import HomeSongCard from './HomeSongCard.tsx'
import Song from '../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'
import { afterEach, expect } from 'vitest'

describe('Home Song Card', () => {
  const song: Song = {
    ...emptySong,
    title: 'Song 1'
  }

  afterEach(() => (window.location.pathname = '/'))

  it('should render with minimal info', () => {
    reduxRouterRender(<HomeSongCard song={song} />)

    expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
    expect(screen.getByText(song.title)).toBeInTheDocument()
  })

  it('should render with maximal info', () => {
    const localSong: Song = {
      ...song,
      imageUrl: 'something.png',
      artist: {
        ...emptyArtist,
        name: 'Artist 1'
      }
    }

    reduxRouterRender(<HomeSongCard song={localSong} />)

    expect(screen.getByRole('img', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.title })).toHaveAttribute(
      'src',
      localSong.imageUrl
    )
    expect(screen.getByText(localSong.title)).toBeInTheDocument()
    expect(screen.getByText(localSong.artist.name)).toBeInTheDocument()
  })

  it('should open song drawer by clicking on image', async () => {
    const user = userEvent.setup()

    const [_, store] = reduxRouterRender(<HomeSongCard song={song} />)

    await user.click(screen.getByLabelText(`default-icon-${song.title}`))
    expect((store.getState() as RootState).global.songDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.songDrawer.songId).toBe(song.id)
  })

  it('should open artist drawer by clicking on artist name', async () => {
    const user = userEvent.setup()

    const localSong: Song = {
      ...song,
      artist: {
        ...emptyArtist,
        name: 'Artist 1'
      }
    }

    const [_, store] = reduxRouterRender(<HomeSongCard song={localSong} />)

    await user.click(screen.getByText(localSong.artist.name))
    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localSong.artist.id)
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRouterRender(<HomeSongCard song={song} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`default-icon-${song.title}`)
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeDisabled()

    rerender(
      <HomeSongCard
        song={{ ...song, artist: emptyArtist, album: emptyAlbum, songsterrLink: 'some-link' }}
      />
    )
    expect(screen.getByRole('menuitem', { name: /view artist/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /open links/i })).not.toBeDisabled()
  })

  describe('on menu', () => {
    it('should navigate to song when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(withToastify(<HomeSongCard song={song} />))

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${song.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/song/${song.id}`)
    })

    it('should navigate to artist when clicking on view artist', async () => {
      const user = userEvent.setup()

      const localSong = { ...song, artist: { ...emptyArtist, id: '1' } }

      reduxRouterRender(<HomeSongCard song={localSong} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${localSong.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view artist/i }))
      expect(window.location.pathname).toBe(`/artist/${localSong.artist.id}`)
    })

    it('should navigate to album when clicking on view album', async () => {
      const user = userEvent.setup()

      const localSong = { ...song, album: { ...emptyAlbum, id: '1' } }

      reduxRouterRender(<HomeSongCard song={localSong} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${localSong.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view album/i }))
      expect(window.location.pathname).toBe(`/album/${localSong.album.id}`)
    })
  })
})
