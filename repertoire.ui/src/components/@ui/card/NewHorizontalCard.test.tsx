import { mantineRender } from '../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { expect } from 'vitest'
import NewHorizontalCard from './NewHorizontalCard.tsx'

describe('New Horizontal Card', () => {
  it('should render and on click event', async () => {
    // Arrange
    const user = userEvent.setup()

    const onClick = vitest.fn()

    const testId = 'children-mock'
    const children = <div data-testid={testId}>Children</div>

    // Act
    mantineRender(<NewHorizontalCard onClick={onClick}>{children}</NewHorizontalCard>)

    // Assert
    expect(screen.getByTestId(testId)).toBeInTheDocument()

    const card = screen.getByLabelText('edit-panel-card')
    await user.click(card)
    expect(onClick).toHaveBeenCalledOnce()
  })
})
