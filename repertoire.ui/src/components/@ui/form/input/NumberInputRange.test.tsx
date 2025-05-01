import { screen } from '@testing-library/react'
import NumberInputRange from './NumberInputRange.tsx'
import { userEvent } from '@testing-library/user-event'
import { mantineRender } from '../../../../test-utils.tsx'

describe('NumberInputRange', () => {
  it('should render inputs ', () => {
    const label = 'Price Range'

    mantineRender(
      <NumberInputRange
        label={label}
        value1={10}
        onChange1={vi.fn()}
        value2={20}
        onChange2={vi.fn()}
      />
    )

    expect(screen.getByText(label)).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: `min ${label}` })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: `min ${label}` })).not.toBeDisabled()
    expect(screen.getByText('-')).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: `max ${label}` })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: `max ${label}` })).not.toBeDisabled()
  })

  it('should call onChange handlers', async () => {
    const user = userEvent.setup()

    const label = 'Price Range'
    const onChange1 = vi.fn()
    const onChange2 = vi.fn()
    const newValue1 = 8
    const newValue2 = 12

    mantineRender(
      <NumberInputRange
        label={label}
        value1={5}
        onChange1={onChange1}
        value2={15}
        onChange2={onChange2}
      />
    )

    await user.clear(screen.getByRole('textbox', { name: `min ${label}` }))
    await user.type(screen.getByRole('textbox', { name: `min ${label}` }), newValue1.toString())
    expect(onChange1).toHaveBeenCalledWith(newValue1)

    await user.clear(screen.getByRole('textbox', { name: `max ${label}` }))
    await user.type(screen.getByRole('textbox', { name: `max ${label}` }), newValue2.toString())
    expect(onChange2).toHaveBeenCalledWith(newValue2)
  })

  it('should disable inputs when loading', () => {
    const label = 'Price Range'
    const onChange1 = vi.fn()
    const onChange2 = vi.fn()

    mantineRender(
      <NumberInputRange
        label={label}
        value1={0}
        onChange1={onChange1}
        value2={100}
        onChange2={onChange2}
        isLoading={true}
      />
    )

    expect(screen.getByRole('textbox', { name: `min ${label}` })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: `max ${label}` })).toBeDisabled()
  })

  it('should enforce max value constraint on first input', async () => {
    const user = userEvent.setup()

    const label = 'Price Range'
    const onChange1 = vi.fn()
    const onChange2 = vi.fn()
    const maxValue = 60
    const newValue = 70

    mantineRender(
      <NumberInputRange
        label={label}
        value1={50}
        onChange1={onChange1}
        value2={100}
        onChange2={onChange2}
        max={maxValue}
      />
    )

    await user.type(screen.getByRole('textbox', { name: `min ${label}` }), newValue.toString())

    expect(onChange1).not.toHaveBeenCalledWith(newValue)
  })

  it('should prevent negative and decimal values', async () => {
    const user = userEvent.setup()

    const label = 'Price Range'
    const onChange1 = vi.fn()
    const onChange2 = vi.fn()
    const newNegativeValue = -12
    const newDecimalValue = 12.5

    mantineRender(
      <NumberInputRange
        label={label}
        value1={10}
        onChange1={onChange1}
        value2={20}
        onChange2={onChange2}
      />
    )

    await user.clear(screen.getByRole('textbox', { name: `min ${label}` }))
    await user.type(screen.getByRole('textbox', { name: `min ${label}` }), newNegativeValue.toString())
    await user.type(screen.getByRole('textbox', { name: `min ${label}` }), newDecimalValue.toString())
    expect(onChange1).not.toHaveBeenCalledWith(newNegativeValue)
    expect(onChange1).not.toHaveBeenCalledWith(newDecimalValue)

    await user.clear(screen.getByRole('textbox', { name: `max ${label}` }))
    await user.type(screen.getByRole('textbox', { name: `max ${label}` }), newNegativeValue.toString())
    await user.type(screen.getByRole('textbox', { name: `max ${label}` }), newDecimalValue.toString())
    expect(onChange2).not.toHaveBeenCalledWith(newNegativeValue)
    expect(onChange2).not.toHaveBeenCalledWith(newDecimalValue)
  })
})
