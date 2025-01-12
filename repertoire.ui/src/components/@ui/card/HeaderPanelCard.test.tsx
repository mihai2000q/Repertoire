import { mantineRender } from '../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { expect } from 'vitest'
import HeaderPanelCard from './HeaderPanelCard.tsx'

describe('Header Panel Card', () => {
  it('should render, display content and buttons on hover', async () => {
    const user = userEvent.setup()

    const onEditClick = vitest.fn()

    const childrenTestId = 'children-mock'
    const children = <div data-testid={childrenTestId}>Children</div>

    const menuTestId = 'menu-mock'
    const menuDropdown = <div data-testid={menuTestId}>This is the menu</div>

    mantineRender(
      <HeaderPanelCard menuDropdown={menuDropdown} onEditClick={onEditClick}>
        {children}
      </HeaderPanelCard>
    )

    expect(screen.getByTestId(childrenTestId)).toBeInTheDocument()

    const moreButton = await screen.findByRole('button', { name: 'more-menu' })
    const editButton = await screen.findByRole('button', { name: 'edit-header' })

    expect(moreButton).not.toBeVisible()
    expect(editButton).not.toBeVisible()
    expect(screen.queryByTestId(menuTestId)).not.toBeInTheDocument()

    const card = screen.getByLabelText('header-panel-card')
    await user.hover(card)

    expect(moreButton).toBeVisible()
    expect(editButton).toBeVisible()

    await user.click(moreButton)
    expect(await screen.findByTestId(menuTestId)).toBeVisible()

    await user.click(editButton)
    expect(onEditClick).toHaveBeenCalledOnce()
  })

  it('should hide icons', async () => {
    mantineRender(
      <HeaderPanelCard menuDropdown={<></>} onEditClick={() => {}} hideIcons={true}>
        Children
      </HeaderPanelCard>
    )

    expect(screen.queryByRole('button', { name: 'more-menu' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'edit-header' })).not.toBeInTheDocument()
  })
})
