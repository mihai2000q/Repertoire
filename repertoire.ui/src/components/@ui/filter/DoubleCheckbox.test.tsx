import { screen } from '@testing-library/react'
import DoubleCheckbox from './DoubleCheckbox.tsx'
import { mantineRender } from '../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'

describe('Double Checkbox', () => {
  it('should render both checkboxes with title, labels and check', () => {
    const title = 'Test title'
    const label1 = 'Option 1'
    const label2 = 'Option 2'

    mantineRender(
      <DoubleCheckbox
        checked1={false}
        onChange1={vi.fn()}
        checked2={false}
        onChange2={vi.fn()}
        title={title}
        label1={label1}
        label2={label2}
      />
    )

    expect(screen.getByText(title)).toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: label1 })).toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: label1 })).not.toBeChecked()
    expect(screen.getByRole('checkbox', { name: label1 })).not.toBeDisabled()
    expect(screen.getByRole('checkbox', { name: label2 })).toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: label2 })).not.toBeDisabled()
  })

  it('should update first checkbox and uncheck second', async () => {
    const user = userEvent.setup()

    const label1 = 'Option 1'
    const onChange1 = vi.fn()
    const onChange2 = vi.fn()

    const { rerender } = mantineRender(
      <DoubleCheckbox
        label1={label1}
        checked1={false}
        onChange1={onChange1}
        checked2={true}
        onChange2={onChange2}
      />
    )

    await user.click(screen.getByRole('checkbox', { name: label1 }))

    expect(onChange1).toHaveBeenCalledExactlyOnceWith(true)

    rerender(
      <DoubleCheckbox
        label1={label1}
        checked1={true}
        onChange1={onChange1}
        checked2={true}
        onChange2={onChange2}
      />
    )
    expect(onChange2).toHaveBeenCalledExactlyOnceWith(false)
  })

  it('should update second checkbox and uncheck first', async () => {
    const user = userEvent.setup()

    const label2 = 'Option 1'
    const onChange1 = vi.fn()
    const onChange2 = vi.fn()

    const { rerender } = mantineRender(
      <DoubleCheckbox
        checked1={true}
        onChange1={onChange1}
        label2={label2}
        checked2={false}
        onChange2={onChange2}
      />
    )

    await user.click(screen.getByRole('checkbox', { name: label2 }))

    expect(onChange2).toHaveBeenCalledExactlyOnceWith(true)

    rerender(
      <DoubleCheckbox
        checked1={true}
        onChange1={onChange1}
        label2={label2}
        checked2={true}
        onChange2={onChange2}
      />
    )
    expect(onChange1).toHaveBeenCalledExactlyOnceWith(false)
  })

  it('should not uncheck the other box when unchecking', async () => {
    const user = userEvent.setup()

    const label1 = 'Option 1'
    const onChange1 = vi.fn()
    const onChange2 = vi.fn()

    mantineRender(
      <DoubleCheckbox
        checked1={true}
        onChange1={onChange1}
        checked2={false}
        onChange2={onChange2}
        label1={label1}
      />
    )

    await user.click(screen.getByRole('checkbox', { name: label1 }))

    expect(onChange1).toHaveBeenCalledExactlyOnceWith(false)
    expect(onChange2).not.toHaveBeenCalled()
  })

  it('should disable both checkboxes when is disabled', () => {
    const label1 = 'Option 1'
    const label2 = 'Option 2'

    mantineRender(
      <DoubleCheckbox
        checked1={false}
        onChange1={vi.fn()}
        checked2={false}
        onChange2={vi.fn()}
        label1={label1}
        label2={label2}
        disabled={true}
      />
    )

    expect(screen.getByRole('checkbox', { name: label1 })).toBeDisabled()
    expect(screen.getByRole('checkbox', { name: label2 })).toBeDisabled()
  })
})
