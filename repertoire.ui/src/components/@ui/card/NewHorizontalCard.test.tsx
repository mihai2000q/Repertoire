import { mantineRender } from '../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { expect } from 'vitest'
import NewHorizontalCard from './NewHorizontalCard.tsx'

describe('New Horizontal Card', () => {
  it('should render and on click event', async () => {
    const user = userEvent.setup()

    const onClick = vitest.fn()

    const testId = 'children-mock'
    const children = <span data-testid={testId}>Children</span>

    mantineRender(<NewHorizontalCard onClick={onClick}>{children}</NewHorizontalCard>)

    expect(screen.getByTestId(testId)).toBeInTheDocument()
    await user.click(screen.getByTestId(testId))
    expect(onClick).toHaveBeenCalledOnce()
  })
})
