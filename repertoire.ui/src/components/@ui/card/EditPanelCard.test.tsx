import { mantineRender } from '../../../test-utils.tsx'
import EditPanelCard from './EditPanelCard.tsx'
import { userEvent } from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { expect } from 'vitest'

describe('Edit Panel Card', () => {
  it('should render and display edit button on hover', async () => {
    // Arrange
    const user = userEvent.setup()

    const onEditClick = vitest.fn()

    const testId = 'children-mock'
    const children = <div data-testid={testId}>Children</div>

    // Act
    mantineRender(<EditPanelCard onEditClick={onEditClick}>{children}</EditPanelCard>)

    // Assert
    expect(screen.getByTestId(testId)).toBeInTheDocument()

    const button = await screen.findByRole('button', { name: 'edit-panel' })

    expect(button).not.toBeVisible()

    const card = screen.getByLabelText('edit-panel-card')
    await user.hover(card)

    expect(button).toBeVisible()

    await user.click(button)
    expect(onEditClick).toHaveBeenCalledOnce()
  })
})
