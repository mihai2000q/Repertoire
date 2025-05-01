import { mantineRender } from '../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import ConfidenceBar from './ConfidenceBar.tsx'

describe('Confidence Bar', () => {
  it('should render', async () => {
    const user = userEvent.setup()

    const confidence = 15

    mantineRender(<ConfidenceBar confidence={confidence} />)

    expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    await user.hover(screen.getByRole('progressbar', { name: 'confidence' }))
    expect(screen.getByRole('tooltip')).toBeInTheDocument()
    expect(screen.getByText(confidence.toString())).toBeInTheDocument()
  })
})
