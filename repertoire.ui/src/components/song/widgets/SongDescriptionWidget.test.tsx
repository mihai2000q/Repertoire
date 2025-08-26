import SongDescriptionWidget from './SongDescriptionWidget.tsx'
import Song from '../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { emptySong, reduxRender } from '../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'

describe('Song Description Widget', () => {
  const song: Song = {
    ...emptySong,
    description: 'This is a description of the song'
  }

  it('should render when there is a description', () => {
    reduxRender(<SongDescriptionWidget song={song} />)

    expect(screen.getByRole('button', { name: 'edit-widget' })).toBeInTheDocument()
    expect(screen.getByText(song.description)).toBeInTheDocument()
  })

  it('should render when there is no description', () => {
    reduxRender(<SongDescriptionWidget song={{ ...song, description: '' }} />)

    expect(screen.getByRole('button', { name: 'edit-widget' })).toBeInTheDocument()
    expect(screen.getByText(/no description/i)).toBeInTheDocument()
  })

  it('should open edit song description modal on edit panel click', async () => {
    const user = userEvent.setup()

    reduxRender(<SongDescriptionWidget song={song} />)

    await user.click(screen.getByRole('button', { name: 'edit-widget' }))
    expect(
      await screen.findByRole('dialog', { name: /edit song description/i })
    ).toBeInTheDocument()
  })
})
