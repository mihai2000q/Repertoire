import {
  ActionIcon,
  ActionIconProps,
  NumberInput,
  NumberInputProps,
  Popover,
  Tooltip
} from '@mantine/core'
import { ReactNode, useState } from 'react'

interface NumberInputButtonProps extends ActionIconProps {
  icon?: ReactNode
  inputProps?: NumberInputProps
  isSelected?: boolean
  inputKey?: string
  tooltipLabels?: {
    selected: string
    default: string
  }
}

function NumberInputButton({
  icon,
  isSelected,
  inputKey,
  inputProps,
  tooltipLabels,
  ...others
}: NumberInputButtonProps) {
  const [opened, setOpened] = useState(false)

  isSelected ??=
    inputProps?.value !== undefined && inputProps?.value !== null && inputProps?.value !== ''

  return (
    <Popover
      opened={opened}
      onChange={setOpened}
      transitionProps={{ transition: 'scale-y', duration: 160 }}
      shadow={'sm'}
      trapFocus
      withArrow
    >
      <Popover.Target>
        <Tooltip
          disabled={opened || tooltipLabels === undefined}
          label={
            inputProps.error ??
            (tooltipLabels &&
              (isSelected === true ? tooltipLabels.selected : tooltipLabels.default))
          }
          openDelay={500}
        >
          <ActionIcon
            variant={'form'}
            aria-selected={isSelected}
            aria-invalid={!!inputProps.error}
            onClick={() => setOpened(!opened)}
            {...others}
          >
            {icon}
          </ActionIcon>
        </Tooltip>
      </Popover.Target>

      <Popover.Dropdown miw={180} p={'xxs'}>
        <NumberInput variant={'unstyled'} size={'xs'} key={inputKey} {...inputProps} />
      </Popover.Dropdown>
    </Popover>
  )
}

export default NumberInputButton
