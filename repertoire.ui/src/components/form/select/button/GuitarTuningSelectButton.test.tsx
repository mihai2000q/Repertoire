import { reduxRender } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import GuitarTuningSelectButton from './GuitarTuningSelectButton.tsx'
import { GuitarTuning } from '../../../../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'

describe('GuitarTuning Select Button', () => {
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

  const handlers = [
    http.get('/songs/guitar-tunings', async () => {
      return HttpResponse.json(guitarTunings)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and change tunings', async () => {
    const user = userEvent.setup()

    const newGuitarTuning = guitarTunings[1]

    const setGuitarTuning = vitest.fn()

    const [{ rerender }] = reduxRender(
      <GuitarTuningSelectButton guitarTuning={null} setGuitarTuning={setGuitarTuning} />
    )

    const button = screen.getByRole('button', { name: 'guitar-tuning' })
    await user.click(button)

    expect(screen.getByRole('textbox', { name: /search/i })).toHaveValue('')
    guitarTunings.forEach((gt) => {
      expect(screen.getByRole('option', { name: gt.name })).toBeInTheDocument()
    })

    const selectedOption = screen.getByRole('option', { name: newGuitarTuning.name })
    await user.click(selectedOption)

    expect(setGuitarTuning).toHaveBeenCalledOnce()
    expect(setGuitarTuning).toHaveBeenCalledWith(newGuitarTuning)

    rerender(
      <GuitarTuningSelectButton guitarTuning={newGuitarTuning} setGuitarTuning={setGuitarTuning} />
    )

    await user.click(button)
    expect(screen.getByRole('textbox', { name: /search/i })).toHaveValue(newGuitarTuning.name)
  })
})
