import dayjs from 'dayjs'
import { Text, TextProps } from '@mantine/core'
import { useRef } from 'react'

interface PastDateProps extends TextProps {
  dateValue: string | null | undefined
  dateFormat?: string
  nullValue?: string
}

function PastDate({ dateValue, dateFormat = 'DD MMM', ...props }: PastDateProps) {
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

  return <Text {...props}>{dateText && dateText}</Text>
}

export default PastDate
