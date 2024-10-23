import { screen } from '@testing-library/react'
import SongCard from './SongCard'
import Song from '../../types/models/Song'
import { mantineRender } from '../../test-utils'

describe('Song Card', () => {
  it('should render and display minimal info', ({ expect }) => {
    // Arrange
    const song: Song = {
      id: '',
      title: 'Some song',
      isRecorded: false
    }

    // Act
    const { container } = mantineRender(<SongCard song={song} />)

    // Assert
    expect(screen.getByText(song.title)).toBeVisible()
    expect(container.querySelector('svg.tabler-icon-microphone-filled')).not.toBeInTheDocument()
  })

  it('should render and display maximal info', ({ expect }) => {
    // Arrange
    const song: Song = {
      id: '',
      title: 'Some song',
      isRecorded: true
    }

    // Act
    const { container } = mantineRender(<SongCard song={song} />)

    // Assert
    expect(screen.getByText(song.title)).toBeVisible()
    expect(container.querySelector('svg.tabler-icon-microphone-filled')).toBeVisible()
  })
})
