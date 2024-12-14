import { screen } from '@testing-library/react'
import SongCard from './SongCard'
import Song from '../../types/models/Song'
import { reduxRouterRender } from '../../test-utils'

describe('Song Card', () => {
  it('should render and display minimal info', ({ expect }) => {
    // Arrange
    const song: Song = {
      id: '',
      title: 'Some song',
      description: '',
      isRecorded: false,
      sections: [],
      rehearsals: 0,
      confidence: 0,
      progress: 0
    }

    // Act
    const [{ container }] = reduxRouterRender(<SongCard song={song} />)

    // Assert
    expect(screen.getByText(song.title)).toBeVisible()
    expect(container.querySelector('svg.tabler-icon-microphone-filled')).not.toBeInTheDocument()
  })

  it('should render and display maximal info', ({ expect }) => {
    // Arrange
    const song: Song = {
      id: '',
      title: 'Some song',
      description: '',
      isRecorded: true,
      sections: [],
      rehearsals: 0,
      confidence: 0,
      progress: 0
    }

    // Act
    const [{ container }] = reduxRouterRender(<SongCard song={song} />)

    // Assert
    expect(screen.getByText(song.title)).toBeVisible()
    expect(container.querySelector('svg.tabler-icon-microphone-filled')).toBeVisible()
  })
})
