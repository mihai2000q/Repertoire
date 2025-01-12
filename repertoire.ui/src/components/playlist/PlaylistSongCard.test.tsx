import { reduxRender } from '../../test-utils.tsx'
import PlaylistSongCard from './PlaylistSongCard.tsx'
import Song from '../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Artist from '../../types/models/Artist.ts'
import Album from '../../types/models/Album.ts'
import { RootState } from '../../state/store.ts'

describe('Playlist Song Card', () => {
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
    playlistTrackNo: 1
  }

  const album: Album = {
    id: '1',
    title: 'Album 1',
    createdAt: '',
    updatedAt: '',
    songs: []
  }

  const artist: Artist = {
    id: '1',
    name: 'Artist 1',
    createdAt: '',
    updatedAt: '',
    albums: [],
    songs: []
  }

  it('should render and display minimal info', () => {
    reduxRender(<PlaylistSongCard song={song} handleRemove={() => {}} />)

    expect(screen.getByText(song.playlistTrackNo)).toBeInTheDocument()
    expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
    expect(screen.getByText(song.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should render and display maximal info', () => {
    const localSong = {
      ...song,
      artist: artist,
      album: album
    }

    reduxRender(<PlaylistSongCard song={localSong} handleRemove={() => {}} />)

    expect(screen.getByText(localSong.playlistTrackNo)).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByText(localSong.title)).toBeInTheDocument()
    expect(screen.getByText(localSong.album.title)).toBeInTheDocument()
    expect(screen.getByText(localSong.artist.name)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    reduxRender(<PlaylistSongCard song={song} handleRemove={() => {}} />)

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display warning modal and remove, when clicking on remove', async () => {
      // Arrange
      const user = userEvent.setup()

      const handleRemove = vitest.fn()

      // Act
      reduxRender(<PlaylistSongCard song={song} handleRemove={handleRemove} />)

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /remove/i }))

      expect(screen.getByRole('heading', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(handleRemove).toHaveBeenCalledOnce()
    })
  })

  it('should open album drawer on album title click', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      album: album
    }

    const [_, store] = reduxRender(<PlaylistSongCard song={localSong} handleRemove={() => {}} />)

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

    const [_, store] = reduxRender(<PlaylistSongCard song={localSong} handleRemove={() => {}} />)

    await user.click(screen.getByText(localSong.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localSong.artist.id)
  })
})
