import { reduxRender } from '../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Album from 'src/types/models/Album.ts'
import ArtistAlbumCard from './ArtistAlbumCard.tsx'

describe('Artist Album Card', () => {
  const album: Album = {
    id: '1',
    title: 'Album 1',
    createdAt: '',
    updatedAt: '',
    songs: []
  }

  it('should render and display minimal information', async () => {
    reduxRender(<ArtistAlbumCard album={album} handleRemove={() => { }} isUnknownArtist={false} />)

    expect(screen.getByRole('img', { name: album.title })).toBeInTheDocument()
    expect(screen.getByText(album.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should render and display maximal information', async () => {
    // Arrange
    const localAlbum: Album = {
      ...album,
      releaseDate: '2024-10-11'
    }

    // Act
    reduxRender(<ArtistAlbumCard album={localAlbum} handleRemove={() => { }} isUnknownArtist={false} />)

    // Assert
    expect(screen.getByRole('img', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByText(localAlbum.title)).toBeInTheDocument()
    expect(screen.getByText('11 Oct 2024')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should display menu by clicking on the dots button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(<ArtistAlbumCard album={album} handleRemove={() => { }} isUnknownArtist={false} />)

    // Assert
    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()
  })

  it('should display less information on the menu when the artist is unknown', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(<ArtistAlbumCard album={album} handleRemove={() => { }} isUnknownArtist={true} />)

    // Assert
    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.queryByRole('menuitem', { name: /remove/i })).not.toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display warning modal and remove album, when clicking on remove', async () => {
      // Arrange
      const user = userEvent.setup()

      const handleRemove = vitest.fn()

      // Act
      reduxRender(<ArtistAlbumCard album={album} handleRemove={handleRemove} isUnknownArtist={false} />)

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /remove/i }))

      expect(screen.getByRole('heading', { name: /remove album/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(handleRemove).toHaveBeenCalledOnce()
    })
  })
})
