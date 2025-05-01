import { DatePicker } from '@mantine/dates'
import { ReactNode, useState } from 'react'
import { ActionIcon, ActionIconProps, alpha, Popover } from '@mantine/core'
import { IconCalendar, IconCalendarCheck } from '@tabler/icons-react'

interface DatePickerButtonProps extends ActionIconProps {
  icon: ReactNode
  value: Date | null
  onChange: (value: Date | null) => void
}

function DatePickerButton({ icon, value, onChange, ...others }: DatePickerButtonProps) {
  const [opened, setOpened] = useState(false)

  return (
    <Popover
      opened={opened}
      onChange={setOpened}
      shadow={'sm'}
      transitionProps={{ transition: 'scale-y', duration: 160 }}
    >
      <Popover.Target>
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
          {value !== null ? <IconCalendarCheck size={20} /> : (icon ?? <IconCalendar size={20} />)}
        </ActionIcon>
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
