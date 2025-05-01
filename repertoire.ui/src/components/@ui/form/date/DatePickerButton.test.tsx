import { screen } from '@testing-library/react'
import DatePickerButton from './DatePickerButton.tsx'
import { IconCalendar } from '@tabler/icons-react'
import { mantineRender } from '../../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'
import dayjs from 'dayjs'

describe('Date Picker Button', () => {
  it('should render', async () => {
    const user = userEvent.setup()

    const ariaLabel = 'Date Picker Button'
    const onChange = vi.fn()

    const now = new Date()
    const newReleaseDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0)

    mantineRender(
      <DatePickerButton
        aria-label={ariaLabel}
        value={null}
        onChange={onChange}
        icon={<IconCalendar data-testid="custom-icon" />}
      />
    )

    expect(screen.getByRole('button', { name: ariaLabel })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: ariaLabel }))
    expect(await screen.findByRole('dialog')).toBeInTheDocument()
    await user.click(
      screen.getByRole('button', { name: dayjs(newReleaseDate).format('D MMMM YYYY') })
    )

    expect(onChange).toHaveBeenCalledExactlyOnceWith(newReleaseDate)
  })
})
