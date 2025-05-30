import { reduxRender } from '../../../../../test-utils.tsx'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { GuitarTuning } from '../../../../../types/models/Song.ts'
import GuitarTuningMultiSelect from './GuitarTuningMultiSelect.tsx'

describe('Guitar Tuning Multi Select', () => {
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
      name: 'Drop C'
    }
  ]

  const handlers = [
    http.get('/songs/guitar-tunings', async () => {
      return HttpResponse.json(guitarTunings)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and change guitar tunings', async () => {
    const user = userEvent.setup()

    const newTunings = [guitarTunings[0], guitarTunings[1]]

    const label = 'guitar-tunings'
    const setIds = vitest.fn()

    reduxRender(<GuitarTuningMultiSelect ids={[]} setIds={setIds} label={label} />)

    const multiSelect = screen.getByRole('textbox', { name: label })
    expect(multiSelect).toHaveValue('')
    expect(multiSelect).toBeDisabled()
    await waitFor(() => expect(multiSelect).not.toBeDisabled())

    await user.click(multiSelect)
    for (const tuning of guitarTunings) {
      expect(await screen.findByRole('option', { name: tuning.name })).toBeInTheDocument()
    }

    for (const tuning of newTunings) {
      await user.click(screen.getByRole('option', { name: tuning.name }))
    }

    expect(setIds).toHaveBeenCalledTimes(newTunings.length)
    newTunings.reduce((a: string[], b) => {
      if (a.length !== 0) expect(setIds).toHaveBeenCalledWith(a) // skip the first case
      return [...a, b.id]
    }, [])
  })
})
