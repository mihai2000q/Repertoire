import { emptyOrder, reduxRouterRender } from '../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Album from 'src/types/models/Album.ts'
import ArtistAlbumCard from './ArtistAlbumCard.tsx'
import dayjs from 'dayjs'
import { expect } from 'vitest'
import AlbumProperty from '../../utils/enums/AlbumProperty.ts'

describe('Artist Album Card', () => {
  const album: Album = {
    id: '1',
    title: 'Album 1',
    createdAt: '',
    updatedAt: '',
    songs: []
  }

  it('should render and display minimal information', async () => {
    reduxRouterRender(
      <ArtistAlbumCard
        album={album}
        handleRemove={() => {}}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    expect(screen.getByRole('img', { name: album.title })).toBeInTheDocument()
    expect(screen.getByText(album.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should render and display maximal information', async () => {
    const localAlbum: Album = {
      ...album,
      imageUrl: 'something.png'
    }

    reduxRouterRender(
      <ArtistAlbumCard
        album={localAlbum}
        handleRemove={() => {}}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    expect(screen.getByRole('img', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.title })).toHaveAttribute(
      'src',
      localAlbum.imageUrl
    )
    expect(screen.getByText(localAlbum.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  describe('on order property change', () => {
    it('should display the release date when it is release date', async () => {
      const localAlbum: Album = {
        ...album,
        releaseDate: '2024-10-11'
      }

      reduxRouterRender(
        <ArtistAlbumCard
          album={localAlbum}
          handleRemove={() => {
          }}
          isUnknownArtist={false}
          order={{...emptyOrder, property: AlbumProperty.ReleaseDate}}
        />
      )

      expect(screen.getByText(dayjs(localAlbum.releaseDate).format('D MMM YYYY'))).toBeInTheDocument()
    })
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistAlbumCard
        album={album}
        handleRemove={() => {}}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`album-card-${album.title}`)
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistAlbumCard
        album={album}
        handleRemove={() => {}}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()
  })

  it('should display less information on the menu when the artist is unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistAlbumCard
        album={album}
        handleRemove={() => {}}
        isUnknownArtist={true}
        order={emptyOrder}
      />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.queryByRole('menuitem', { name: /remove/i })).not.toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to album when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistAlbumCard
          album={album}
          handleRemove={() => {}}
          isUnknownArtist={false}
          order={emptyOrder}
        />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))

      expect(window.location.pathname).toBe(`/album/${album.id}`)

      // restore
      window.location.pathname = '/'
    })

    it('should display warning modal and remove album, when clicking on remove', async () => {
      const user = userEvent.setup()

      const handleRemove = vitest.fn()

      reduxRouterRender(
        <ArtistAlbumCard
          album={album}
          handleRemove={handleRemove}
          isUnknownArtist={false}
          order={emptyOrder}
        />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /remove/i }))

      expect(await screen.findByRole('dialog', { name: /remove album/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /remove album/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(handleRemove).toHaveBeenCalledOnce()
    })
  })
})
