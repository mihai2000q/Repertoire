import { screen } from '@testing-library/react'
import DatePickerButton from './DatePickerButton.tsx'
import { mantineRender } from '../../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'
import dayjs from 'dayjs'

describe('Date Picker Button', () => {
  it('should render', async () => {
    const user = userEvent.setup()

    const ariaLabel = 'Date Picker Button'
    const onChange = vi.fn()

    const newDate = dayjs().format('YYYY-MM-DD')

    const { rerender } = mantineRender(
      <DatePickerButton aria-label={ariaLabel} value={null} onChange={onChange} />
    )

    const button = screen.getByRole('button', { name: ariaLabel })

    expect(button).toBeInTheDocument()

    await user.hover(button)
    expect(await screen.findByRole('tooltip', { name: 'Select a date' })).toBeInTheDocument()

    await user.click(button)
    expect(await screen.findByRole('dialog')).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: dayjs(newDate).format('D MMMM YYYY') }))

    expect(onChange).toHaveBeenCalledExactlyOnceWith(newDate)

    rerender(<DatePickerButton aria-label={ariaLabel} value={newDate} onChange={onChange} />)

    await user.hover(button)
    expect(
      await screen.findByRole('tooltip', {
        name: `${dayjs(newDate).format('D MMMM YYYY')} is selected`
      })
    ).toBeInTheDocument()
  })
})
