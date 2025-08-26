import { mantineRender } from '../../../test-utils.tsx'
import EditWidget from './EditWidget.tsx'
import { userEvent } from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { expect } from 'vitest'

describe('Edit Widget', () => {
  it('should render and display edit button on hover', async () => {
    const user = userEvent.setup()

    const onEditClick = vitest.fn()

    const testId = 'children-mock'
    const children = <div data-testid={testId}>Children</div>

    mantineRender(<EditWidget onEditClick={onEditClick}>{children}</EditWidget>)

    expect(screen.getByTestId(testId)).toBeInTheDocument()

    const button = await screen.findByRole('button', { name: 'edit-widget' })

    expect(button).not.toBeVisible()

    await user.hover(screen.getByTestId(testId))

    expect(button).toBeVisible()

    await user.hover(button)
    expect(await screen.findByText(/edit widget/i)).toBeInTheDocument()

    await user.click(button)
    expect(onEditClick).toHaveBeenCalledOnce()
  })
})
