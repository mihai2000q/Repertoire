import { mantineRender } from '../../../../test-utils.tsx'
import DifficultySelect from './DifficultySelect.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Difficulty from '../../../../utils/enums/Difficulty.ts'
import { expect } from 'vitest'

describe('Difficulty Select', () => {
  it('should render and change difficulties', async () => {
    // Arrange
    const user = userEvent.setup()

    const [difficultyLabel, difficultyValue] = Object.entries(Difficulty)[1]
    const newOption = { label: difficultyLabel, value: difficultyValue }

    const onChange = vitest.fn()

    // Act
    const { rerender } = mantineRender(<DifficultySelect option={null} onChange={onChange} />)

    // Assert
    expect(screen.getByRole('textbox', { name: /difficulty/i })).toHaveValue('')

    const select = screen.getByRole('textbox', { name: /difficulty/i })
    await user.click(select)

    Object.entries(Difficulty).forEach(([key]) => {
      expect(screen.getByRole('option', { name: key })).toBeInTheDocument()
    })

    const selectedOption = screen.getByRole('option', { name: difficultyLabel })
    await user.click(selectedOption)

    expect(onChange).toHaveBeenCalledOnce()
    expect(onChange).toHaveBeenCalledWith(newOption)

    rerender(<DifficultySelect option={newOption} onChange={onChange} />)

    expect(screen.getByRole('textbox', { name: /difficulty/i })).toHaveValue(difficultyLabel)
  })
})
