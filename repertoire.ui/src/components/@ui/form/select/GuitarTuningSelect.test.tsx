import { reduxRender } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
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
    // Arrange
    const user = userEvent.setup()

    const guitarTuning = guitarTunings[0]
    const newOption = { label: guitarTuning.name, value: guitarTuning.id }

    const onChange = vitest.fn()

    // Act
    const [{ rerender }] = reduxRender(<GuitarTuningSelect option={null} onChange={onChange} />)

    // Assert
    expect(screen.getByText(/loading/i)).toBeInTheDocument()

    expect(await screen.findByRole('textbox', { name: /guitar tuning/i })).toHaveValue('')

    const select = screen.getByRole('textbox', { name: /guitar tuning/i })
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
