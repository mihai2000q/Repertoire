import { reduxRender } from '../../test-utils.tsx'
import ArtistSongCard from './ArtistSongCard.tsx'
import Song from '../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Album from 'src/types/models/Album.ts'

describe('Artist Song Card', () => {
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
  }

  const album: Album = {
    id: '1',
    title: 'Album 1',
    createdAt: '',
    updatedAt: '',
    songs: []
  }

  it('should render and display minimal information', async () => {
    reduxRender(<ArtistSongCard song={song} handleRemove={() => { }} isUnknownArtist={false} />)

    expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
    expect(screen.getByText(song.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should render and display maximal information', async () => {
    // Arrange
    const user = userEvent.setup()

    const localSong: Song = {
      ...song,
      releaseDate: '2024-10-11',
      album: album
    }

    // Act
    reduxRender(<ArtistSongCard song={localSong} handleRemove={() => { }} isUnknownArtist={false} />)

    // Assert
    expect(screen.getByRole('img', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByText(localSong.title)).toBeInTheDocument()
    expect(screen.getByText('11 Oct 2024')).toBeInTheDocument()
    expect(screen.getByText(localSong.album.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should display menu by clicking on the dots button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(<ArtistSongCard song={song} handleRemove={() => { }} isUnknownArtist={false} />)

    // Assert
    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()
  })

  it('should display less information on the menu when the artist is unknown', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(<ArtistSongCard song={song} handleRemove={() => { }} isUnknownArtist={true} />)

    // Assert
    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.queryByRole('menuitem', { name: /remove/i })).not.toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display warning modal and remove song, when clicking on remove', async () => {
      // Arrange
      const user = userEvent.setup()

      const handleRemove = vitest.fn()

      // Act
      reduxRender(<ArtistSongCard song={song} handleRemove={handleRemove} isUnknownArtist={false} />)

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /remove/i }))

      expect(screen.getByRole('heading', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(handleRemove).toHaveBeenCalledOnce()
    })
  })
})
