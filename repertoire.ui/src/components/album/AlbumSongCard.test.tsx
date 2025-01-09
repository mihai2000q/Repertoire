import { reduxRender } from '../../test-utils.tsx'
import AlbumSongCard from './AlbumSongCard.tsx'
import Song from '../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Album Song Card', () => {
  const song: Song = {
    id: '1',
    title: 'Song 1',
    description: '',
    isRecorded: false,
    rehearsals: 0,
    confidence: 0,
    progress: 0,
    sections: [],
    createdAt: '',
    updatedAt: '',
    albumTrackNo: 1
  }

  it('should render and display information, when the album is not unknown', () => {
    reduxRender(<AlbumSongCard song={song} handleRemove={() => { }} isUnknownAlbum={false} />)

    expect(screen.getByText(song.albumTrackNo)).toBeInTheDocument()
    expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
    expect(screen.getByText(song.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should display menu by clicking on the dots button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(<AlbumSongCard song={song} handleRemove={() => { }} isUnknownAlbum={false} />)

    // Assert
    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display warning modal and remove, when clicking on remove', async () => {
      // Arrange
      const user = userEvent.setup()

      const handleRemove = vitest.fn()

      // Act
      reduxRender(<AlbumSongCard song={song} handleRemove={handleRemove} isUnknownAlbum={false} />)

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /remove/i }))

      expect(screen.getByRole('heading', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(handleRemove).toHaveBeenCalledOnce()
    })
  })

  it('should not display the tracking number and some menu options, when the album is unknown', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(<AlbumSongCard song={song} handleRemove={() => { }} isUnknownAlbum={true} />)

    // Assert
    expect(screen.queryByText(song.albumTrackNo)).not.toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.queryByRole('menuitem', { name: /remove/i })).not.toBeInTheDocument()
  })
})
