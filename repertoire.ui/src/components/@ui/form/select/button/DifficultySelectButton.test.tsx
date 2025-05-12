import { mantineRender } from '../../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Difficulty from '../../../../../types/enums/Difficulty.ts'
import { expect } from 'vitest'
import DifficultySelectButton from './DifficultySelectButton.tsx'

describe('Difficulty Select Button', () => {
  it('should render and change difficulties', async () => {
    const user = userEvent.setup()

    const [difficultyLabel, difficultyValue] = Object.entries(Difficulty)[1]
    const newDifficulty = difficultyValue

    const setDifficulty = vitest.fn()

    const { rerender } = mantineRender(
      <DifficultySelectButton difficulty={null} setDifficulty={setDifficulty} />
    )

    const button = screen.getByRole('button', { name: 'difficulty' })
    await user.click(button)

    expect(screen.getByRole('textbox', { name: /search/i })).toHaveValue('')
    Object.entries(Difficulty).forEach(([_, value]) => {
      expect(screen.getByRole('option', { name: value })).toBeInTheDocument()
    })

    const selectedOption = screen.getByRole('option', { name: difficultyValue })
    await user.click(selectedOption)

    expect(setDifficulty).toHaveBeenCalledOnce()
    expect(setDifficulty).toHaveBeenCalledWith(newDifficulty)

    rerender(<DifficultySelectButton difficulty={newDifficulty} setDifficulty={setDifficulty} />)

    await user.click(button)
    expect(screen.getByRole('textbox', { name: /search/i })).toHaveValue(difficultyLabel)
  })
})
