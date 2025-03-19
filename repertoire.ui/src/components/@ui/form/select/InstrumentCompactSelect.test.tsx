import { reduxRender } from '../../../../test-utils.tsx'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { Instrument } from '../../../../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import InstrumentCompactSelect from './InstrumentCompactSelect.tsx'

describe('Instrument Compact Select', () => {
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

    const onChange = vitest.fn()

    const [{ rerender }] = reduxRender(
      <InstrumentCompactSelect instrument={null} setInstrument={onChange} />
    )

    const selectButton = screen.getByRole('button', { name: 'select-instrument' })
    expect(selectButton).toBeDisabled()
    await waitFor(() => expect(selectButton).not.toBeDisabled())
    await user.click(selectButton)

    for (const instrument of instruments) {
      expect(await screen.findByRole('instrument', { name: instrument.name })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('instrument', { name: newInstrument.name })
    await user.click(selectedOption)

    expect(onChange).toHaveBeenCalledOnce()
    expect(onChange).toHaveBeenCalledWith(newInstrument)

    rerender(<InstrumentCompactSelect instrument={newInstrument} setInstrument={onChange} />)

    expect(screen.queryByRole('button', { name: 'select-band-member' })).not.toBeInTheDocument()

    const instrumentButton = screen.getByRole('button', { name: newInstrument.name })
    expect(instrumentButton).toBeInTheDocument()

    await user.hover(instrumentButton)
    expect(
      await screen.findByRole('tooltip', { name: new RegExp(newInstrument.name, 'i') })
    ).toBeInTheDocument()

    await user.click(instrumentButton)
    expect(screen.getByRole('textbox', { name: /search/i })).toHaveValue(newInstrument.name)

    // reset the value from outside component
    rerender(<InstrumentCompactSelect instrument={null} setInstrument={onChange} />)

    expect(screen.getByRole('button', { name: 'select-instrument' })).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: 'select-instrument' }))
    expect(screen.getByRole('textbox', { name: /search/i })).toHaveValue('')
  })

  it('should filter by name', async () => {
    const user = userEvent.setup()

    const searchValue = 't'

    reduxRender(<InstrumentCompactSelect instrument={null} setInstrument={() => {}} />)

    await user.click(screen.getByRole('button', { name: 'select-instrument' }))
    await user.type(screen.getByRole('textbox', { name: /search/i }), searchValue)

    const filteredInstruments = instruments.filter((i) => i.name.includes(searchValue))
    expect(await screen.findAllByRole('instrument')).toHaveLength(filteredInstruments.length)
    for (const instrument of filteredInstruments) {
      expect(screen.getByRole('instrument', { name: instrument.name })).toBeInTheDocument()
    }
  })
})
