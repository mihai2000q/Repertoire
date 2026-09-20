import { mantineRender } from '../../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import SongOutlineConfidenceBar from './SongOutlineConfidenceBar.tsx'

describe('Song Outline Confidence Bar', () => {
  it('should render', async () => {
    const user = userEvent.setup()

    const confidence = 65

    mantineRender(<SongOutlineConfidenceBar confidence={confidence} />)

    expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    await user.hover(screen.getByRole('progressbar', { name: 'confidence' }))
    expect(screen.getByRole('tooltip')).toBeInTheDocument()
    expect(screen.getByText(`Confidence: ${confidence}%`)).toBeInTheDocument()
  })
})
