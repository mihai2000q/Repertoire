import { userEvent } from '@testing-library/user-event'
import { mantineRender } from '../../../test-utils.tsx'
import { ContextMenu } from './ContextMenu.tsx'
import { screen } from '@testing-library/react'

describe('Context Menu', () => {
  it('should render and display menu on right click', async () => {
    const user = userEvent.setup()

    const buttonName = 'click me'
    const dropdownTestId = 'dropdownTestId'
    const menuItem = 'contextMenuItem'

    mantineRender(
      <ContextMenu>
        <ContextMenu.Target>
          <button>{buttonName}</button>
        </ContextMenu.Target>

        <ContextMenu.Dropdown datatest-id={dropdownTestId}>
          <ContextMenu.Item>{menuItem}</ContextMenu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>
    )

    expect(screen.queryByTestId(dropdownTestId)).not.toBeInTheDocument()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByRole('button', { name: buttonName })
    })

    expect(await screen.findByTestId(dropdownTestId)).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: menuItem })).toBeInTheDocument()
  })
})
