import { mantineRender } from '../../../test-utils.tsx'
import EditPanelCard from './EditPanelCard.tsx'
import { userEvent } from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { expect } from 'vitest'

describe('Edit Panel Card', () => {
  it('should render and display edit button on hover', async () => {
    const user = userEvent.setup()

    const onEditClick = vitest.fn()

    const testId = 'children-mock'
    const children = <div data-testid={testId}>Children</div>

    mantineRender(<EditPanelCard onEditClick={onEditClick}>{children}</EditPanelCard>)

    expect(screen.getByTestId(testId)).toBeInTheDocument()

    const button = await screen.findByRole('button', { name: 'edit-panel' })

    expect(button).not.toBeVisible()

    await user.hover(screen.getByTestId(testId))

    expect(button).toBeVisible()

    await user.click(button)
    expect(onEditClick).toHaveBeenCalledOnce()
  })
})
