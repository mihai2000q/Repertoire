import { mantineRender } from '../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import ProgressBar from './ProgressBar.tsx'

describe('Progress Bar', () => {
  it('should render', async () => {
    const user = userEvent.setup()

    const progress = 15

    mantineRender(<ProgressBar progress={progress} />)

    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    await user.hover(screen.getByRole('progressbar', { name: 'progress' }))
    expect(screen.getByRole('tooltip')).toBeInTheDocument()
    expect(screen.getByText(progress.toString())).toBeInTheDocument()
  })
})
