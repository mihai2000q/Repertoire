import Song from '../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { mantineRender } from '../../../test-utils.tsx'
import SongOverallCard from './SongOverallCard.tsx'
import {userEvent} from "@testing-library/user-event";

describe('Song Overall Card', () => {
  const song: Song = {
    id: '',
    title: '',
    description: '',
    isRecorded: false,
    sections: [],
    rehearsals: 10,
    confidence: 78,
    progress: 254,
    createdAt: '',
    updatedAt: ''
  }

  it('should render', async () => {
    const user = userEvent.setup()

    mantineRender(<SongOverallCard song={song} />)

    expect(screen.getByText(song.rehearsals)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toHaveValue(song.confidence)
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toHaveValue(song.progress / 10)

    await user.hover(screen.getByRole('progressbar', { name: 'confidence' }))
    expect(
      screen.getByRole('tooltip', { name: new RegExp(song.confidence.toString()) })
    ).toBeInTheDocument()

    await user.hover(screen.getByRole('progressbar', { name: 'progress' }))
    expect(
      screen.getByRole('tooltip', { name: new RegExp(song.progress.toString()) })
    ).toBeInTheDocument()
  })
})
