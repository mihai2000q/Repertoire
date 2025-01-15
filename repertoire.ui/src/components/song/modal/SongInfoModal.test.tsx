import { mantineRender } from '../../../test-utils.tsx'
import SongInfoModal from './SongInfoModal.tsx'
import Song from '../../../types/models/Song.ts'
import { setupServer } from 'msw/node'
import { screen } from '@testing-library/react'
import dayjs from 'dayjs'

describe('Song Info Modal', () => {
  const song: Song = {
    id: '',
    title: '',
    description: '',
    isRecorded: false,
    sections: [],
    rehearsals: 0,
    confidence: 0,
    progress: 0,
    createdAt: '2024-10-15T10:30',
    updatedAt: '2024-10-16T22:16'
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    mantineRender(<SongInfoModal opened={true} onClose={() => {}} song={song} />)

    expect(screen.getByRole('dialog', { name: /song info/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /song info/i })).toBeInTheDocument()
    expect(
      screen.getByText(dayjs(song.createdAt).format('DD MMMM YYYY, HH:mm'))
    ).toBeInTheDocument()
    expect(
      screen.getByText(dayjs(song.updatedAt).format('DD MMMM YYYY, HH:mm'))
    ).toBeInTheDocument()
  })
})
