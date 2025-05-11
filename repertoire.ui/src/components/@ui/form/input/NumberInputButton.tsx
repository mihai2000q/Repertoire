import {
  ActionIcon,
  ActionIconProps,
  NumberInput,
  NumberInputProps,
  Popover,
  Tooltip
} from '@mantine/core'
import { forwardRef, ReactNode, useState } from 'react'

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

const NumberInputButton = forwardRef<HTMLButtonElement, NumberInputButtonProps>(
  ({ icon, isSelected, inputKey, inputProps, tooltipLabels, ...others }, ref) => {
    const [opened, setOpened] = useState(false)

    isSelected ??=
      inputProps?.value !== undefined && inputProps?.value !== null && inputProps?.value !== ''

    return (
      <Popover
        opened={opened}
        onChange={setOpened}
        transitionProps={{ transition: 'scale-y', duration: 160 }}
        trapFocus
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
              ref={ref}
              size={'lg'}
              variant={'form'}
              aria-selected={isSelected}
              onClick={() => setOpened(!opened)}
              {...others}
              c={inputProps.error && 'red.7'}
              bg={inputProps.error && 'red.2'}
              styles={(theme) => ({
                root: { ...(inputProps.error && { border: `2px solid ${theme.colors.red[5]}` }) }
              })}
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
)

NumberInputButton.displayName = 'NumberInputButton'

export default NumberInputButton
