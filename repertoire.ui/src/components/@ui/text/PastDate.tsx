import dayjs from 'dayjs'
import { Text, TextProps } from '@mantine/core'
import { forwardRef, useRef } from 'react'

interface PastDateProps extends TextProps {
  dateValue: string | null | undefined
  dateFormat?: string
  nullValue?: string
}

const PastDate = forwardRef<HTMLDivElement, PastDateProps>(
  ({ dateValue, dateFormat = 'DD MMM', ...props }, ref) => {
    const now = useRef(dayjs())

    const date = dateValue !== null && dateValue !== undefined ? dayjs(dateValue) : undefined
    const dateText = date
      ? date.isToday()
        ? 'Today'
        : date.isYesterday()
          ? 'Yesterday'
          : date.isSame(now.current, 'week')
            ? date.format('dddd')
            : date.format(dateFormat)
      : undefined

    return (
      <Text ref={ref} {...props}>
        {dateText && dateText}
      </Text>
    )
  }
)

PastDate.displayName = 'PastDate'

export default PastDate
