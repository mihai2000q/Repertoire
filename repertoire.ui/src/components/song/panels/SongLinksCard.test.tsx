import Song from '../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { reduxRender } from '../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'
import SongLinksCard from './SongLinksCard.tsx'

describe('Song Links Card', () => {
  const song: Song = {
    id: '',
    title: '',
    description: '',
    isRecorded: false,
    sections: [],
    rehearsals: 0,
    confidence: 0,
    progress: 0,
    createdAt: '',
    updatedAt: '',
    youtubeLink: 'https://www.youtube.com/watch?v=123456789',
    songsterrLink: 'https://www.songsterr.com/song'
  }

  it('should render when there are links', () => {
    reduxRender(<SongLinksCard song={song} />)

    expect(screen.getByRole('button', { name: 'edit-panel' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /youtube/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /songsterr/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /youtube/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /songsterr/i })).toBeInTheDocument()
  })

  it('should render when there are no links', () => {
    reduxRender(<SongLinksCard song={{ ...song, youtubeLink: '', songsterrLink: '' }} />)

    expect(screen.getByRole('button', { name: 'edit-panel' })).toBeInTheDocument()
    expect(screen.getByText(/no links/i)).toBeInTheDocument()
  })

  it('should open edit song links modal on edit panel click', async () => {
    const user = userEvent.setup()

    reduxRender(<SongLinksCard song={song} />)

    await user.click(screen.getByRole('button', { name: 'edit-panel' }))

    expect(await screen.findByRole('dialog', { name: /edit song links/i })).toBeInTheDocument()
  })
})
