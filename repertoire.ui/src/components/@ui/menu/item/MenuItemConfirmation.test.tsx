import { mantineRender } from '../../../../test-utils.tsx'
import { Menu } from '@mantine/core'
import { expect } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import MenuItemConfirmation from './MenuItemConfirmation.tsx'
import { ReactNode } from 'react'

describe('Menu Item Confirmation', () => {
  const render = (children: ReactNode) =>
    mantineRender(
      <Menu opened={true}>
        <Menu.Dropdown>{children}</Menu.Dropdown>
      </Menu>
    )

  it('should render', () => {
    render(
      <MenuItemConfirmation onConfirm={vi.fn()}>
        <span data-testid={'test-item'}>Something</span>
      </MenuItemConfirmation>
    )

    expect(screen.getByTestId('test-item')).toBeInTheDocument()
  })

  it('should display cancel and confirm buttons and invoke methods', async () => {
    const user = userEvent.setup()

    const onCancel = vi.fn()
    const onConfirm = vi.fn()

    const name = 'something'

    render(
      <MenuItemConfirmation onConfirm={onConfirm} onCancel={onCancel}>
        {name}
      </MenuItemConfirmation>
    )

    await user.click(screen.getByRole('menuitem', { name: name }))

    expect(screen.getByRole('button', { name: 'cancel' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'confirm' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'cancel' }))
    expect(onCancel).toHaveBeenCalledOnce()
    await waitFor(() => {
      expect(screen.queryByRole('button', { name: 'cancel' })).not.toBeInTheDocument()
      expect(screen.queryByRole('button', { name: 'confirm' })).not.toBeInTheDocument()
    })

    await user.click(screen.getByRole('menuitem', { name: name }))

    await user.click(screen.getByRole('button', { name: 'confirm' }))
    expect(onConfirm).toHaveBeenCalledOnce()
    await waitFor(() => {
      expect(screen.queryByRole('button', { name: 'cancel' })).not.toBeInTheDocument()
      expect(screen.queryByRole('button', { name: 'confirm' })).not.toBeInTheDocument()
    })
  })
})
