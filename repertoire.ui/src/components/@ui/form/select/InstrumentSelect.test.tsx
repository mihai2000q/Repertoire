import { reduxRender } from '../../../../test-utils.tsx'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { Instrument } from '../../../../types/models/Song.ts'
import { setupServer } from 'msw/node'
import InstrumentSelect from './InstrumentSelect.tsx'

describe('Instrument Select', () => {
  const instruments: Instrument[] = [
    {
      id: '1',
      name: 'Guitar'
    },
    {
      id: '2',
      name: 'Piano'
    },
    {
      id: '3',
      name: 'Flute'
    }
  ]

  const handlers = [
    http.get('/songs/instruments', async () => {
      return HttpResponse.json(instruments)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and change instruments', async () => {
    const user = userEvent.setup()

    const newInstrument = instruments[0]
    const newOption = { label: newInstrument.name, value: newInstrument.id }

    const onChange = vitest.fn()

    const label = 'label'

    const [{ rerender }] = reduxRender(
      <InstrumentSelect label={label} option={null} onOptionChange={onChange} />
    )

    const select = screen.getByRole('textbox', { name: label })
    expect(select).toHaveValue('')
    expect(select).toBeDisabled()
    await waitFor(() => expect(select).not.toBeDisabled())
    await user.click(select)

    for (const instrument of instruments) {
      expect(await screen.findByRole('option', { name: instrument.name })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('option', { name: newInstrument.name })
    await user.click(selectedOption)

    expect(onChange).toHaveBeenCalledOnce()
    expect(onChange).toHaveBeenCalledWith(newOption)

    rerender(<InstrumentSelect label={label} option={newOption} onOptionChange={onChange} />)

    expect(screen.getByRole('textbox', { name: label })).toHaveValue(newInstrument.name)
  })
})
