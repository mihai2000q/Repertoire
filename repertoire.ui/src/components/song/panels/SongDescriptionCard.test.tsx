import SongDescriptionCard from './SongDescriptionCard.tsx'
import Song from '../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { reduxRender } from '../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'

describe('Song Description Card', () => {
  const song: Song = {
    id: '',
    title: '',
    description: 'This is a description of the song',
    isRecorded: false,
    sections: [],
    rehearsals: 0,
    confidence: 0,
    progress: 0,
    createdAt: '',
    updatedAt: ''
  }

  it('should render when there is a description', () => {
    reduxRender(<SongDescriptionCard song={song} />)

    expect(screen.getByRole('button', { name: 'edit-panel' })).toBeInTheDocument()
    expect(screen.getByText(song.description)).toBeInTheDocument()
  })

  it('should render when there is no description', () => {
    reduxRender(<SongDescriptionCard song={{ ...song, description: '' }} />)

    expect(screen.getByRole('button', { name: 'edit-panel' })).toBeInTheDocument()
    expect(screen.getByText(/no description/i)).toBeInTheDocument()
  })

  it('should open edit song description modal on edit panel click', async () => {
    const user = userEvent.setup()

    reduxRender(<SongDescriptionCard song={song} />)

    await user.click(screen.getByRole('button', { name: 'edit-panel' }))
    expect(screen.getByRole('dialog', { name: /edit song description/i })).toBeInTheDocument()
  })
})
