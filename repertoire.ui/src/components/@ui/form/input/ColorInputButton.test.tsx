import { mantineRender } from '../../../../test-utils.tsx'
import ColorInputButton from './ColorInputButton.tsx'
import {screen} from "@testing-library/react";
import {userEvent} from "@testing-library/user-event";

describe('Color Input Button', () => {
  it('should render', async () => {
    const user = userEvent.setup()

    const { rerender } = mantineRender(<ColorInputButton color={''} setColor={() => {}} />)

    expect(screen.getByText(/color/i)).toBeInTheDocument() // label
    const colorInput = screen.getByRole('button', { name: 'color-input' })
    expect(colorInput).toBeInTheDocument()
    expect(colorInput).toHaveStyle('backgroundColor: \'transparent\'')

    await user.click(colorInput)
    expect(screen.getByRole('dialog')).toBeInTheDocument()

    const newColor = '#123'
    rerender(<ColorInputButton color={newColor} setColor={() => {}} />)

    expect(screen.getByRole('button', { name: 'color-input' })).toHaveStyle(`backgroundColor: '${newColor}'`)
  })
})
