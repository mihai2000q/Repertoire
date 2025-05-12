import { ActionIcon, alpha, Group, Popover, PopoverProps, Text } from '@mantine/core'
import { IconCheck, IconX } from '@tabler/icons-react'
import { ReactNode } from 'react'

interface PopoverConfirmationProps {
  children: ReactNode
  label: string
  popoverProps: PopoverProps
  isLoading?: boolean
  onCancel?: () => void
  onConfirm?: () => void
}

function PopoverConfirmation({
  children,
  label,
  popoverProps,
  isLoading,
  onCancel,
  onConfirm
}: PopoverConfirmationProps) {
  return (
    <Popover transitionProps={{ transition: 'fade-up' }} position={'top'} {...popoverProps}>
      <Popover.Target>{children}</Popover.Target>

      <Popover.Dropdown>
        <Group gap={'xxs'}>
          <Text c={'dimmed'} fw={500} fz={'sm'}>
            {label}
          </Text>
          <Group gap={'xxs'}>
            <ActionIcon
              variant={'subtle'}
              aria-label={'cancel'}
              disabled={isLoading}
              onClick={onCancel}
              sx={(theme) => ({
                color: theme.colors.red[4],
                '&:hover': {
                  color: theme.colors.red[5],
                  backgroundColor: alpha(theme.colors.red[2], 0.35)
                },
                '&[data-disabled]': {
                  color: theme.colors.gray[4],
                  backgroundColor: 'transparent'
                }
              })}
            >
              <IconX size={16} />
            </ActionIcon>
            <ActionIcon
              variant={'subtle'}
              c={'green'}
              aria-label={'confirm'}
              loading={isLoading}
              onClick={onConfirm}
              sx={(theme) => ({
                '&:hover': { backgroundColor: alpha(theme.colors.green[2], 0.35) }
              })}
            >
              <IconCheck size={16} />
            </ActionIcon>
          </Group>
        </Group>
      </Popover.Dropdown>
    </Popover>
  )
}

export default PopoverConfirmation
