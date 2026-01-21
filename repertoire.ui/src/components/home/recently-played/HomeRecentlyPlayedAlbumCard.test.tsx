import { emptyAlbum, emptyArtist, reduxRouterRender } from '../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { afterEach, expect } from 'vitest'
import { RootState } from '../../../state/store.ts'
import Album from '../../../types/models/Album.ts'
import HomeRecentlyPlayedAlbumCard from './HomeRecentlyPlayedAlbumCard.tsx'
import dayjs from 'dayjs'
import { userEvent } from '@testing-library/user-event'

describe('Home Recently Played Album Card', () => {
  const album: Album = {
    ...emptyAlbum,
    id: '1',
    imageUrl: 'something.png',
    title: 'Album 1',
    lastTimePlayed: '2024-12-14T10:30',
    progress: 500,
    artist: {
      ...emptyArtist,
      id: '1',
      name: 'another-artist'
    }
  }

  afterEach(() => (window.location.pathname = '/'))

  it('should render', () => {
    const localAlbum = {
      ...album,
      imageUrl: null,
      artist: null
    }

    reduxRouterRender(<HomeRecentlyPlayedAlbumCard album={localAlbum} />)

    expect(screen.getByText(localAlbum.title)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByLabelText(`default-icon-${localAlbum.title}`)).toBeInTheDocument()
    expect(screen.getByText(dayjs(localAlbum.lastTimePlayed).format('DD MMM'))).toBeInTheDocument()
  })

  it('should render with artist', () => {
    reduxRouterRender(<HomeRecentlyPlayedAlbumCard album={album} />)

    expect(screen.getByText(album.artist.name)).toBeInTheDocument()
  })

  it('should render with image', () => {
    reduxRouterRender(<HomeRecentlyPlayedAlbumCard album={album} />)

    expect(screen.getByRole('img', { name: album.title })).toHaveAttribute('src', album.imageUrl)
  })

  it('should open album drawer on clicking', async () => {
    const user = userEvent.setup()

    const [_, store] = reduxRouterRender(<HomeRecentlyPlayedAlbumCard album={album} />)

    await user.click(await screen.findByText(album.title))

    expect((store.getState() as RootState).global.albumDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.albumDrawer.albumId).toBe(album.id)
  })

  it('should open artist drawer on clicking the artist', async () => {
    const user = userEvent.setup()

    const [_, store] = reduxRouterRender(<HomeRecentlyPlayedAlbumCard album={album} />)

    await user.click(await screen.findByText(album.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(album.artist.id)
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    const localAlbum: Album = {
      ...album,
      artist: null
    }

    const [{ rerender }] = reduxRouterRender(<HomeRecentlyPlayedAlbumCard album={localAlbum} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: await screen.findByRole('img', { name: localAlbum.title })
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeDisabled()

    rerender(<HomeRecentlyPlayedAlbumCard album={album} />)
    await user.pointer({
      keys: '[MouseRight>]',
      target: await screen.findByRole('img', { name: album.title })
    })
    expect(screen.getByRole('menuitem', { name: /view artist/i })).not.toBeDisabled()
  })

  describe('on menu', () => {
    it('should navigate to album when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<HomeRecentlyPlayedAlbumCard album={album} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: await screen.findByRole('img', { name: album.title })
      })
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/album/${album.id}`)
    })

    it('should navigate to artist when clicking on view artist', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<HomeRecentlyPlayedAlbumCard album={album} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: await screen.findByRole('img', { name: album.title })
      })
      await user.click(screen.getByRole('menuitem', { name: /view artist/i }))
      expect(window.location.pathname).toBe(`/artist/${album.artist.id}`)
    })
  })
})
