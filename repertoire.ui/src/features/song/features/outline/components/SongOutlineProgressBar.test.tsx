import { mantineRender } from '../../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import SongOutlineProgressBar from './SongOutlineProgressBar.tsx'

describe('Song Outline Progress Bar', () => {
  it('should render', async () => {
    const user = userEvent.setup()

    const progress = 15
    const maxProgress = 30

    mantineRender(<SongOutlineProgressBar progress={progress} maxProgress={maxProgress} />)

    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    await user.hover(screen.getByRole('progressbar', { name: 'progress' }))
    expect(screen.getByRole('tooltip')).toBeInTheDocument()
    expect(screen.getByText('Progress:')).toBeInTheDocument()
    expect(screen.getByText(progress.toString())).toBeInTheDocument()
  })

  it('should compute value as a percentage of maxProgress', () => {
    mantineRender(<SongOutlineProgressBar progress={15} maxProgress={30} />)

    expect(screen.getByRole('progressbar')).toHaveAttribute('aria-valuenow', '50')
  })

  it('should render 0 when progress is 0', () => {
    mantineRender(<SongOutlineProgressBar progress={0} maxProgress={30} />)

    expect(screen.getByRole('progressbar')).toHaveAttribute('aria-valuenow', '0')
  })
})
