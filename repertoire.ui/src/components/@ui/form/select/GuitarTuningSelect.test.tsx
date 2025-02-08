import { reduxRender } from '../../../../test-utils.tsx'
import {screen, waitFor} from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { GuitarTuning } from '../../../../types/models/Song.ts'
import { setupServer } from 'msw/node'
import GuitarTuningSelect from './GuitarTuningSelect.tsx'

describe('Guitar Tuning Select', () => {
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

    const guitarTuning = guitarTunings[0]
    const newOption = { label: guitarTuning.name, value: guitarTuning.id }

    const onChange = vitest.fn()

    const [{ rerender }] = reduxRender(<GuitarTuningSelect option={null} onChange={onChange} />)

    const select = screen.getByRole('textbox', { name: /guitar tuning/i })
    expect(select).toHaveValue('')
    expect(select).toBeDisabled()
    await waitFor(() => expect(select).not.toBeDisabled())
    await user.click(select)

    for (const gt of guitarTunings) {
      expect(await screen.findByRole('option', { name: gt.name })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('option', { name: guitarTuning.name })
    await user.click(selectedOption)

    expect(onChange).toHaveBeenCalledOnce()
    expect(onChange).toHaveBeenCalledWith(newOption)

    rerender(<GuitarTuningSelect option={newOption} onChange={onChange} />)

    expect(screen.getByRole('textbox', { name: /guitar tuning/i })).toHaveValue(guitarTuning.name)
  })
})
