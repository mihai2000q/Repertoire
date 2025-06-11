import { emptyAlbum, emptyArtist, reduxRouterRender } from '../../../test-utils.tsx'
import HomeAlbumCard from './HomeAlbumCard.tsx'
import Album from '../../../types/models/Album.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'
import { afterEach } from 'vitest'

describe('Home Album Card', () => {
  const album: Album = {
    ...emptyAlbum,
    id: '1',
    title: 'Album 1'
  }

  afterEach(() => (window.location.pathname = '/'))

  it('should render with minimal info', () => {
    reduxRouterRender(<HomeAlbumCard album={album} />)

    expect(screen.getByLabelText(`default-icon-${album.title}`)).toBeInTheDocument()
    expect(screen.getByText(album.title)).toBeInTheDocument()
  })

  it('should render with maximal info', () => {
    const localAlbum: Album = {
      ...album,
      imageUrl: 'something.png',
      artist: {
        ...emptyArtist,
        name: 'Artist 1'
      }
    }

    reduxRouterRender(<HomeAlbumCard album={localAlbum} />)

    expect(screen.getByRole('img', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.title })).toHaveAttribute(
      'src',
      localAlbum.imageUrl
    )
    expect(screen.getByText(localAlbum.title)).toBeInTheDocument()
    expect(screen.getByText(localAlbum.artist.name)).toBeInTheDocument()
  })

  it('should open album drawer by clicking on image', async () => {
    const user = userEvent.setup()

    const [_, store] = reduxRouterRender(<HomeAlbumCard album={album} />)

    await user.click(screen.getByLabelText(`default-icon-${album.title}`))
    expect((store.getState() as RootState).global.albumDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.albumDrawer.albumId).toBe(album.id)
  })

  it('should open artist drawer by clicking on artist name', async () => {
    const user = userEvent.setup()

    const localAlbum: Album = {
      ...album,
      artist: {
        ...emptyArtist,
        name: 'Artist 1'
      }
    }

    const [_, store] = reduxRouterRender(<HomeAlbumCard album={localAlbum} />)

    await user.click(screen.getByText(localAlbum.artist.name))
    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localAlbum.artist.id)
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRouterRender(<HomeAlbumCard album={album} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`default-icon-${album.title}`)
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeDisabled()

    rerender(<HomeAlbumCard album={{ ...album, artist: emptyArtist }} />)
    expect(screen.getByRole('menuitem', { name: /view artist/i })).not.toBeDisabled()
  })

  describe('on menu', () => {
    it('should navigate to album when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<HomeAlbumCard album={album} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${album.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/album/${album.id}`)
    })

    it('should navigate to artist when clicking on view artist', async () => {
      const user = userEvent.setup()

      const localAlbum = { ...album, artist: { ...emptyArtist, id: '1' } }

      reduxRouterRender(<HomeAlbumCard album={localAlbum} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${localAlbum.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view artist/i }))
      expect(window.location.pathname).toBe(`/artist/${localAlbum.artist.id}`)
    })
  })
})
