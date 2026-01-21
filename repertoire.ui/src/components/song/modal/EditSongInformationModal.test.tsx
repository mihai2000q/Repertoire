import { emptySong, reduxRender, withToastify } from '../../../test-utils.tsx'
import Song, { GuitarTuning } from '../../../types/models/Song.ts'
import { setupServer } from 'msw/node'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { UpdateSongRequest } from '../../../types/requests/SongRequests.ts'
import EditSongInformationModal from './EditSongInformationModal.tsx'
import Difficulty from '../../../types/enums/Difficulty.ts'

describe('Edit Song Information Modal', () => {
  const guitarTunings: GuitarTuning[] = [
    {
      id: '1',
      name: 'E Standard'
    },
    {
      id: '2',
      name: 'Drop D'
    },
    {
      id: '3',
      name: 'Drop A'
    }
  ]

  const song: Song = {
    ...emptySong,
    id: 'some-id',
    guitarTuning: guitarTunings[1],
    difficulty: Difficulty.Easy,
    bpm: 120,
    isRecorded: true
  }

  const handlers = [
    http.get(`/songs/guitar-tunings`, () => {
      return HttpResponse.json(guitarTunings)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongInformationModal opened={true} onClose={() => {}} song={song} />)

    expect(screen.getByRole('dialog', { name: /edit song information/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit song information/i })).toBeInTheDocument()

    expect(await screen.findByRole('textbox', { name: /guitar tuning/i })).toBeInTheDocument()
    expect(await screen.findByRole('textbox', { name: /guitar tuning/i })).toHaveValue(
      song.guitarTuning.name
    )
    expect(screen.getByRole('textbox', { name: /difficulty/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /difficulty/i })).toHaveValue(
      Difficulty[song.difficulty]
    )
    expect(screen.getByRole('textbox', { name: /bpm/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /bpm/i })).toHaveValue(song.bpm.toString())
    expect(screen.getByRole('checkbox', { name: /recorded/i })).toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: /recorded/i })).toBeChecked()

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send update request when all the information has changed', async () => {
    const user = userEvent.setup()

    const newGuitarTuning = guitarTunings[0]
    const newDifficulty = Difficulty.Hard
    const newBpm = 56
    const onClose = vitest.fn()

    let capturedRequest: UpdateSongRequest
    server.use(
      http.put('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const [{ rerender }] = reduxRender(
      withToastify(<EditSongInformationModal opened={true} onClose={onClose} song={song} />)
    )

    const guitarTuningField = await screen.findByRole('textbox', { name: /guitar tuning/i })
    const difficultyField = screen.getByRole('textbox', { name: /difficulty/i })
    const bpmField = screen.getByRole('textbox', { name: /bpm/i })
    const recordedField = screen.getByRole('checkbox', { name: /recorded/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.click(guitarTuningField)
    await user.click(await screen.findByRole('option', { name: newGuitarTuning.name }))

    await user.click(difficultyField)
    await user.click(await screen.findByText(new RegExp(newDifficulty, 'i')))

    await user.clear(bpmField)
    await user.type(bpmField, newBpm.toString())

    await user.click(recordedField)

    expect(saveButton).not.toHaveAttribute('data-disabled')
    await user.click(saveButton)

    expect(capturedRequest).toStrictEqual({
      ...song,
      guitarTuningId: newGuitarTuning.id,
      difficulty: newDifficulty,
      bpm: newBpm,
      isRecorded: false
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(await screen.findByText(/song information updated/i)).toBeInTheDocument()

    rerender(
      <EditSongInformationModal
        opened={true}
        onClose={onClose}
        song={{
          ...song,
          guitarTuning: newGuitarTuning,
          difficulty: newDifficulty,
          bpm: newBpm,
          isRecorded: false
        }}
      />
    )
    expect(await screen.findByRole('button', { name: /save/i })).toHaveAttribute(
      'data-disabled',
      'true'
    )
  })

  it('should send update request when all the information has changed to empty values', async () => {
    const user = userEvent.setup()

    let capturedRequest: UpdateSongRequest
    server.use(
      http.put('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(<EditSongInformationModal opened={true} onClose={() => {}} song={song} />)

    await user.click(await screen.findByRole('textbox', { name: /guitar tuning/i }))
    await user.click(await screen.findByRole('option', { name: song.guitarTuning.name }))
    await user.click(screen.getByRole('textbox', { name: /difficulty/i }))
    await user.click(await screen.findByText(new RegExp(song.difficulty, 'i')))
    await user.clear(screen.getByRole('textbox', { name: /bpm/i }))
    await user.click(screen.getByRole('checkbox', { name: /recorded/i }))
    await user.click(screen.getByRole('button', { name: /save/i }))

    expect(capturedRequest).toStrictEqual({
      ...song,
      difficulty: null,
      bpm: null,
      isRecorded: false
    })
  })

  it('should keep the save button disabled when the information has not changed', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongInformationModal opened={true} onClose={() => {}} song={song} />)

    const guitarTuningField = await screen.findByRole('textbox', { name: /guitar tuning/i })
    const difficultyField = screen.getByRole('textbox', { name: /difficulty/i })
    const bpmField = screen.getByRole('textbox', { name: /bpm/i })
    const recordedField = screen.getByRole('checkbox', { name: /recorded/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    // clear guitar tuning
    await user.click(guitarTuningField)
    await user.click(await screen.findByRole('option', { name: song.guitarTuning.name }))
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reselect guitar tuning
    await user.click(guitarTuningField)
    await user.click(await screen.findByRole('option', { name: song.guitarTuning.name }))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // clear difficulty
    await user.click(difficultyField)
    await user.click(await screen.findByText(new RegExp(song.difficulty, 'i')))
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reselect difficulty
    await user.click(difficultyField)
    await user.click(await screen.findByText(new RegExp(song.difficulty, 'i')))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change bpm
    await user.clear(bpmField)
    await user.type(bpmField, '123')
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset bpm
    await user.clear(bpmField)
    await user.type(bpmField, song.bpm.toString())
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change recorded
    await user.click(recordedField)
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset recorded
    await user.click(recordedField)
    expect(saveButton).toHaveAttribute('data-disabled', 'true')
  })
})
