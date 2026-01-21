import { mantineRender } from '../../../test-utils.tsx'
import PopoverConfirmation from './PopoverConfirmation.tsx'
import { Button } from '@mantine/core'
import { screen } from '@testing-library/react'
import { expect } from 'vitest'
import { userEvent } from '@testing-library/user-event'

describe('Popover Confirmation', () => {
  it('should render and display popover', async () => {
    const user = userEvent.setup()

    const label = 'this is a label'
    const childrenTestId = 'childrenTestId'
    const children = <Button data-testid={childrenTestId}>Button</Button>

    const { rerender } = mantineRender(
      <PopoverConfirmation label={label} popoverProps={{ opened: false }}>
        {children}
      </PopoverConfirmation>
    )

    expect(screen.getByTestId(childrenTestId)).toBeInTheDocument()
    expect(screen.queryByText(label)).not.toBeInTheDocument()

    const onCancel = vi.fn()
    const onConfirm = vi.fn()

    rerender(
      <PopoverConfirmation
        label={label}
        popoverProps={{ opened: true }}
        onCancel={onCancel}
        onConfirm={onConfirm}
      >
        {children}
      </PopoverConfirmation>
    )

    expect(await screen.findByRole('dialog')).toBeInTheDocument()
    expect(screen.getByText(label)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'cancel' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'confirm' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'cancel' }))
    expect(onCancel).toHaveBeenCalledOnce()

    await user.click(screen.getByRole('button', { name: 'confirm' }))
    expect(onConfirm).toHaveBeenCalledOnce()
  })
})
