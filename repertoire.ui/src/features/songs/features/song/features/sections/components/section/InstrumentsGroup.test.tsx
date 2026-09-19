import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import InstrumentsGroup from './InstrumentsGroup.tsx'
import { Instrument } from '../../../../../../../../types/models/Song.ts'
import { mantineRender } from '../../../../../../../../test-utils.tsx'

describe('Instruments Group', () => {
  const instruments: Instrument[] = [
    { id: '1', name: 'Guitar' },
    { id: '2', name: 'Piano' }
  ]

  it('should render an avatar per instrument', () => {
    mantineRender(<InstrumentsGroup instruments={instruments} />)

    instruments.forEach((instrument) => {
      expect(screen.getByLabelText(instrument.name)).toBeInTheDocument()
    })
  })

  it('should show all instrument names in the tooltip on hover', async () => {
    const user = userEvent.setup()

    mantineRender(<InstrumentsGroup instruments={instruments} />)

    await user.hover(screen.getByLabelText(instruments[0].name))

    for (const instrument of instruments) {
      expect(await screen.findByRole('tooltip', { name: instrument.name })).toBeInTheDocument()
    }
  })
})
