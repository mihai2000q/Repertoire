import { screen } from '@testing-library/react'
import SongCard from './SongCard'
import Song from '../../types/models/Song'
import { reduxRender } from '../../test-utils'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../state/store.ts'

describe('Song Card', () => {
  it('should render and display minimal info', ({ expect }) => {
    // Arrange
    const song: Song = {
      id: '',
      title: 'Some song',
      description: '',
      isRecorded: false,
      sections: []
    }

    // Act
    const [{ container }] = reduxRender(<SongCard song={song} />)

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
      sections: []
    }

    // Act
    const [{ container }] = reduxRender(<SongCard song={song} />)

    // Assert
    expect(screen.getByText(song.title)).toBeVisible()
    expect(container.querySelector('svg.tabler-icon-microphone-filled')).toBeVisible()
  })

  it('should open drawer and dispatch song Id on click', async ({ expect }) => {
    // Arrange
    const song: Song = {
      id: '1',
      title: 'Some song',
      description: '',
      isRecorded: true,
      sections: []
    }

    const user = userEvent.setup()

    // Act
    const [_, store] = reduxRender(<SongCard song={song} />)

    // Assert
    await user.click(screen.getByTestId(`song-card-${song.id}`))

    expect((store.getState() as RootState).global.songDrawer.songId).toBe(song.id)
    expect((store.getState() as RootState).global.songDrawer.open).toBeTruthy()
  })
})
