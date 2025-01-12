import Song from '../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { reduxRender } from '../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'
import SongInformationCard from './SongInformationCard.tsx'
import Difficulty from '../../../utils/enums/Difficulty.ts'
import dayjs from 'dayjs'
import {http, HttpResponse} from "msw";
import {setupServer} from "msw/node";

describe('Song Information Card', () => {
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
    updatedAt: ''
  }

  const handlers = [
    http.get(`/songs/guitar-tunings`, () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render when there is no information', () => {
    reduxRender(<SongInformationCard song={song} />)

    expect(screen.getByRole('button', { name: 'edit-panel' })).toBeInTheDocument()
    expect(screen.getAllByText(/not set/i)).toHaveLength(3)
    expect(screen.getByText('No')).toBeInTheDocument()
    expect(screen.getByText(/never/i)).toBeInTheDocument()
  })

  it('should render when there is information', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      difficulty: Difficulty.Impossible,
      guitarTuning: {
        id: '',
        name: 'Drop D'
      },
      bpm: 120,
      isRecorded: true,
      lastTimePlayed: '2024-10-30'
    }

    reduxRender(<SongInformationCard song={localSong} />)

    expect(screen.getByRole('button', { name: 'edit-panel' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()
    expect(screen.getByText(localSong.guitarTuning.name)).toBeInTheDocument()
    expect(screen.getByText(localSong.bpm)).toBeInTheDocument()
    expect(screen.getByLabelText('recorded-icon')).toBeInTheDocument()
    expect(
      screen.getByText(dayjs(localSong.lastTimePlayed).format('DD MMM YYYY'))
    ).toBeInTheDocument()

    await user.hover(screen.getByRole('progressbar', { name: 'difficulty' }))
    expect(
      await screen.findByRole('tooltip', { name: new RegExp(localSong.difficulty) })
    ).toBeInTheDocument()
  })

  it('should open edit song information modal on edit panel click', async () => {
    const user = userEvent.setup()

    reduxRender(<SongInformationCard song={song} />)

    await user.click(screen.getByRole('button', { name: 'edit-panel' }))
    expect(screen.getByRole('dialog', { name: /edit song information/i })).toBeInTheDocument()
  })
})
