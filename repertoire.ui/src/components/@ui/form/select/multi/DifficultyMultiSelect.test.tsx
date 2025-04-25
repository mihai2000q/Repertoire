import { reduxRender } from '../../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Difficulty from '../../../../../types/enums/Difficulty.ts'
import DifficultyMultiSelect from './DifficultyMultiSelect.tsx'

describe('Difficulty Multi Select', () => {
  const difficulties = Object.entries(Difficulty).map(([_, value]) => value)

  it('should render and change difficulties', async () => {
    const user = userEvent.setup()

    const newDifficulties = [difficulties[0], difficulties[1]]

    const label = 'difficulties'
    const setIds = vitest.fn()

    reduxRender(<DifficultyMultiSelect difficulties={[]} setDifficulties={setIds} label={label} />)

    const multiSelect = screen.getByRole('textbox', { name: label })
    expect(multiSelect).toHaveValue('')

    await user.click(multiSelect)
    for (const difficulty of difficulties) {
      expect(await screen.findByRole('option', { name: difficulty })).toBeInTheDocument()
    }

    for (const difficulty of newDifficulties) {
      await user.click(screen.getByRole('option', { name: difficulty }))
    }

    expect(setIds).toHaveBeenCalledTimes(newDifficulties.length)
    newDifficulties.reduce((a: string[], b) => {
      if (a.length !== 0) expect(setIds).toHaveBeenCalledWith(a) // skip the first case
      return [...a, b]
    }, [])
  })
})
