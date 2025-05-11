import { DatePicker } from '@mantine/dates'
import { ReactNode, useState } from 'react'
import { ActionIcon, ActionIconProps, alpha, Popover, Tooltip } from '@mantine/core'
import { IconCalendar, IconCalendarCheck } from '@tabler/icons-react'
import dayjs from 'dayjs'

interface DatePickerButtonProps extends ActionIconProps {
  value: Date | null | undefined
  onChange: (value: Date | null | undefined) => void
  icon?: ReactNode
  tooltipLabels?: {
    default?: string
    selected?: (date: Date) => string
  }
}

function DatePickerButton({
  icon,
  value,
  onChange,
  tooltipLabels,
  ...others
}: DatePickerButtonProps) {
  const [opened, setOpened] = useState(false)

  return (
    <Popover
      opened={opened}
      onChange={setOpened}
      transitionProps={{ transition: 'scale-y', duration: 160 }}
    >
      <Popover.Target>
        <Tooltip
          disabled={opened}
          label={
            value !== null
              ? tooltipLabels?.selected
                ? tooltipLabels?.selected(value)
                : `${dayjs(value).format('D MMMM YYYY')} is selected`
              : (tooltipLabels?.default ?? 'Select a date')
          }
          openDelay={500}
        >
          <ActionIcon
            variant={value !== null ? 'transparent' : 'grey'}
            size={'lg'}
            sx={(theme) => ({
              ...(value !== null && {
                color: theme.colors.green[5],
                backgroundColor: alpha(theme.colors.green[1], 0.5),

                '&:hover': {
                  color: theme.colors.green[6],
                  backgroundColor: theme.colors.green[1]
                }
              })
            })}
            onClick={() => setOpened(!opened)}
            {...others}
          >
            {value !== null ? (
              <IconCalendarCheck size={20} />
            ) : (
              (icon ?? <IconCalendar size={20} />)
            )}
          </ActionIcon>
        </Tooltip>
      </Popover.Target>

      <Popover.Dropdown>
        <DatePicker
          allowDeselect
          value={value}
          onChange={(val) => {
            onChange(val)
            setOpened(false)
          }}
        />
      </Popover.Dropdown>
    </Popover>
  )
}

export default DatePickerButton
