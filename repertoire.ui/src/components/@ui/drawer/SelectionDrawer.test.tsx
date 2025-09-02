import { mantineRender } from '../../../test-utils.tsx'
import SelectionDrawer from './SelectionDrawer.tsx'
import { screen } from '@testing-library/react'
import { Menu } from '@mantine/core'
import { userEvent } from '@testing-library/user-event'

describe('Selection Drawer', () => {
  it('should render', async () => {
    const text = 'some text'
    const actionIconsTestId = 'actionIcons'

    const { rerender } = mantineRender(
      <SelectionDrawer
        opened={false}
        text={text}
        actionIcons={<div data-testid={actionIconsTestId}></div>}
        onClose={vi.fn()}
      />
    )

    expect(screen.queryByText(text)).not.toBeInTheDocument()

    // open the drawer
    rerender(
      <SelectionDrawer
        opened={true}
        text={text}
        actionIcons={<div data-testid={actionIconsTestId}></div>}
        onClose={vi.fn()}
      />
    )

    expect(await screen.findByText(text)).toBeInTheDocument()
    expect(screen.getByTestId(actionIconsTestId)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'close-drawer' })).toBeInTheDocument()
  })

  it('should render with menu', async () => {
    const user = userEvent.setup()

    const text = 'some text'
    const actionIconsTestId = 'actionIcons'
    const toggleMenu = vi.fn()
    const menuDropdownTestId = 'menuDropdown'

    const { rerender } = mantineRender(
      <SelectionDrawer
        opened={true}
        text={text}
        actionIcons={<div data-testid={actionIconsTestId}></div>}
        onClose={vi.fn()}
        menu={{
          opened: false,
          toggle: toggleMenu,
          dropdown: <Menu.Dropdown data-testid={menuDropdownTestId}></Menu.Dropdown>
        }}
      />
    )

    expect(screen.getByText(text)).toBeInTheDocument()
    expect(screen.getByTestId(actionIconsTestId)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'close-drawer' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(toggleMenu).toHaveBeenCalledExactlyOnceWith(true)

    // open menu
    rerender(
      <SelectionDrawer
        opened={true}
        text={text}
        actionIcons={<div data-testid={actionIconsTestId}></div>}
        onClose={vi.fn()}
        menu={{
          opened: true,
          toggle: toggleMenu,
          dropdown: <Menu.Dropdown data-testid={menuDropdownTestId}></Menu.Dropdown>
        }}
      />
    )

    expect(await screen.findByTestId(menuDropdownTestId)).toBeInTheDocument()
  })

  it('should call onClose when clicking the close button', async () => {
    const user = userEvent.setup()

    const onClose = vi.fn()

    mantineRender(
      <SelectionDrawer
        opened={true}
        text={''}
        actionIcons={<div></div>}
        onClose={onClose}
      />
    )

    await user.click(screen.getByRole('button', { name: 'close-drawer' }))

    expect(onClose).toHaveBeenCalledOnce()
  })
})
