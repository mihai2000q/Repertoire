import { reduxRender } from '../../../../../test-utils.tsx'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { Instrument } from '../../../../../types/models/Song.ts'
import InstrumentMultiSelect from './InstrumentMultiSelect.tsx'

describe('Guitar Tuning Multi Select', () => {
  const instruments: Instrument[] = [
    {
      id: '1',
      name: 'Guitar'
    },
    {
      id: '2',
      name: 'Violin'
    },
    {
      id: '3',
      name: 'Ukulele'
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

  it('should render and change roles', async () => {
    const user = userEvent.setup()

    const newInstruments = [instruments[0], instruments[1]]

    const label = 'instruments'
    const setIds = vitest.fn()

    reduxRender(<InstrumentMultiSelect ids={[]} setIds={setIds} label={label} />)

    const multiSelect = screen.getByRole('textbox', { name: label })
    expect(multiSelect).toHaveValue('')
    expect(multiSelect).toBeDisabled()
    await waitFor(() => expect(multiSelect).not.toBeDisabled())

    await user.click(multiSelect)
    for (const instrument of instruments) {
      expect(await screen.findByRole('option', { name: instrument.name })).toBeInTheDocument()
    }

    for (const instrument of newInstruments) {
      await user.click(screen.getByRole('option', { name: instrument.name }))
    }

    expect(setIds).toHaveBeenCalledTimes(newInstruments.length)
    newInstruments.reduce((a: string[], b) => {
      if (a.length !== 0) expect(setIds).toHaveBeenCalledWith(a) // skip the first case
      return [...a, b.id]
    }, [])
  })

  it('should render and change instruments', async () => {
    const user = userEvent.setup()

    const newInstruments = [instruments[0], instruments[1]]

    const label = 'instruments'
    const setIds = vitest.fn()

    reduxRender(<InstrumentMultiSelect ids={[]} setIds={setIds} label={label} />)

    const multiSelect = screen.getByRole('textbox', { name: label })
    expect(multiSelect).toHaveValue('')
    expect(multiSelect).toBeDisabled()
    await waitFor(() => expect(multiSelect).not.toBeDisabled())

    await user.click(multiSelect)
    for (const instrument of instruments) {
      expect(await screen.findByRole('option', { name: instrument.name })).toBeInTheDocument()
    }

    for (const instrument of newInstruments) {
      await user.click(screen.getByRole('option', { name: instrument.name }))
    }

    expect(setIds).toHaveBeenCalledTimes(newInstruments.length)
    newInstruments.reduce((a: string[], b) => {
      if (a.length !== 0) expect(setIds).toHaveBeenCalledWith(a) // skip the first case
      return [...a, b.id]
    }, [])
  })
})
