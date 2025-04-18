import { mantineRender } from '../../../test-utils.tsx'
import DifficultyBar from './DifficultyBar.tsx'
import Difficulty from '../../../types/enums/Difficulty.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Difficulty Bar', () => {
  const difficulties = Object.entries(Difficulty).map((d) => d[1])

  it.each([...difficulties])('should show difficulty bars', async (difficulty) => {
    mantineRender(<DifficultyBar difficulty={difficulty} />)

    let progressBars = screen.queryAllByRole('progressbar')
    expect(progressBars).toHaveLength(difficulties.length + 1)

    const parentProgressBar = screen.getByRole('progressbar', { name: 'difficulty' })
    progressBars = progressBars.filter((p) => p !== parentProgressBar)
    expect(progressBars).toHaveLength(difficulties.length)

    const difficultyIndex = difficulties.indexOf(difficulty)
    progressBars.slice(0, difficultyIndex).forEach((p) => expect(p).toHaveValue(100))
    progressBars
      .slice(difficultyIndex + 1, progressBars.length)
      .forEach((p) => expect(p).toHaveValue(0))
  })

  it.each([...difficulties])('should show tooltip', async (difficulty) => {
    const user = userEvent.setup()

    mantineRender(<DifficultyBar difficulty={difficulty} />)

    await user.hover(screen.getByRole('progressbar', { name: 'difficulty' }))
    expect(await screen.findByRole('tooltip', { name: new RegExp(difficulty) })).toBeInTheDocument()
  })
})
