import {
  ActionIcon,
  ActionIconProps,
  Popover,
  TextInput,
  TextInputProps,
  Tooltip
} from '@mantine/core'
import { forwardRef, ReactNode, useState } from 'react'

interface TextInputButtonProps extends ActionIconProps {
  icon?: ReactNode
  inputKey?: string
  isSelected?: boolean
  inputProps?: TextInputProps
  tooltipLabels?: {
    selected: string
    default: string
  }
}

const TextInputButton = forwardRef<HTMLButtonElement, TextInputButtonProps>(
  ({ icon, inputKey, isSelected, inputProps, tooltipLabels, ...others }, ref) => {
    const [opened, setOpened] = useState(false)

    isSelected ??=
      inputProps.value !== undefined && inputProps.value !== null && inputProps.value !== ''

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
              variant={'form'}
              aria-selected={isSelected === true}
              aria-invalid={!!inputProps.error}
              onClick={() => setOpened(!opened)}
              {...others}
            >
              {icon}
            </ActionIcon>
          </Tooltip>
        </Popover.Target>

        <Popover.Dropdown miw={180} p={'xxs'}>
          <TextInput variant={'unstyled'} size={'xs'} key={inputKey} {...inputProps} />
        </Popover.Dropdown>
      </Popover>
    )
  }
)

TextInputButton.displayName = 'TextInputButton'

export default TextInputButton
