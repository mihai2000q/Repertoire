import { userEvent } from '@testing-library/user-event'
import { mantineRender } from '../../../test-utils.tsx'
import { ContextMenu } from './ContextMenu.tsx'
import { screen } from '@testing-library/react'

describe('Context Menu', () => {
  it('should render and display menu on right click - uncontrolled', async () => {
    const user = userEvent.setup()

    const buttonName = 'click me'
    const menuItem = 'contextMenuItem'

    mantineRender(
      <ContextMenu>
        <ContextMenu.Target>
          <button>{buttonName}</button>
        </ContextMenu.Target>

        <ContextMenu.Dropdown>
          <ContextMenu.Item>{menuItem}</ContextMenu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>
    )

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByRole('button', { name: buttonName })
    })

    expect(await screen.findByRole('menu')).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: menuItem })).toBeInTheDocument()
  })

  it('should render and display menu on right click - controlled', async () => {
    const user = userEvent.setup()

    const onOpen = vi.fn()
    const onClose = vi.fn()

    const buttonName = 'click me'

    mantineRender(
      <ContextMenu opened={false} onOpen={onOpen} onClose={onClose}>
        <ContextMenu.Target>
          <button>{buttonName}</button>
        </ContextMenu.Target>

        <ContextMenu.Dropdown>
        </ContextMenu.Dropdown>
      </ContextMenu>
    )

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByRole('button', { name: buttonName })
    })

    expect(onOpen).toHaveBeenCalledOnce()
  })

  it('should not open on disabled - uncontrolled', async () => {
    const user = userEvent.setup()

    const buttonName = 'click me'

    mantineRender(
      <ContextMenu disabled={true}>
        <ContextMenu.Target>
          <button>{buttonName}</button>
        </ContextMenu.Target>

        <ContextMenu.Dropdown>
        </ContextMenu.Dropdown>
      </ContextMenu>
    )

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByRole('button', { name: buttonName })
    })

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('should not open on disabled - controlled', async () => {
    const user = userEvent.setup()

    const onOpen = vi.fn()
    const onClose = vi.fn()

    const buttonName = 'click me'

    mantineRender(
      <ContextMenu opened={false} onOpen={onOpen} onClose={onClose} disabled={true}>
        <ContextMenu.Target>
          <button>{buttonName}</button>
        </ContextMenu.Target>

        <ContextMenu.Dropdown>
        </ContextMenu.Dropdown>
      </ContextMenu>
    )

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByRole('button', { name: buttonName })
    })

    expect(onOpen).not.toHaveBeenCalled()
  })
})
