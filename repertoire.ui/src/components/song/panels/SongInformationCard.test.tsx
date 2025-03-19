import { screen } from '@testing-library/react'
import { emptySong, reduxRender } from '../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'
import SongInformationCard from './SongInformationCard.tsx'
import Difficulty from '../../../utils/enums/Difficulty.ts'
import dayjs from 'dayjs'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'

describe('Song Information Card', () => {
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
    reduxRender(<SongInformationCard song={emptySong} />)

    expect(screen.getByRole('button', { name: 'edit-panel' })).toBeInTheDocument()
    expect(screen.getAllByText(/not set/i)).toHaveLength(3)
    expect(screen.getByText('No')).toBeInTheDocument()
    expect(screen.getByText(/never/i)).toBeInTheDocument()
  })

  it('should render when there is information', async () => {
    const user = userEvent.setup()

    const song = {
      ...emptySong,
      difficulty: Difficulty.Impossible,
      guitarTuning: {
        id: '',
        name: 'Drop D'
      },
      bpm: 120,
      isRecorded: true,
      lastTimePlayed: '2024-10-30'
    }

    reduxRender(<SongInformationCard song={song} />)

    expect(screen.getByRole('button', { name: 'edit-panel' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()
    expect(screen.getByText(song.guitarTuning.name)).toBeInTheDocument()
    expect(screen.getByText(song.bpm)).toBeInTheDocument()
    expect(screen.getByLabelText('recorded-icon')).toBeInTheDocument()
    expect(screen.getByText(dayjs(song.lastTimePlayed).format('DD MMM YYYY'))).toBeInTheDocument()

    await user.hover(screen.getByRole('progressbar', { name: 'difficulty' }))
    expect(
      await screen.findByRole('tooltip', { name: new RegExp(song.difficulty) })
    ).toBeInTheDocument()
  })

  it('should open edit song information modal on edit panel click', async () => {
    const user = userEvent.setup()

    reduxRender(<SongInformationCard song={emptySong} />)

    await user.click(screen.getByRole('button', { name: 'edit-panel' }))
    expect(
      await screen.findByRole('dialog', { name: /edit song information/i })
    ).toBeInTheDocument()
  })
})
