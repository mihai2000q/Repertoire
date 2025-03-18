import { IconCheck, IconX } from '@tabler/icons-react'
import { ActionIcon, alpha, Box, Group, Menu, MenuItemProps, Transition } from '@mantine/core'
import { MouseEvent, ReactNode, useState } from 'react'

interface MenuItemConfirmationProps extends MenuItemProps {
  children: ReactNode
  onConfirm?: () => void
  isLoading?: boolean
  onCancel?: () => void
}

function MenuItemConfirmation({
  children,
  isLoading,
  onConfirm,
  onCancel,
  ...props
}: MenuItemConfirmationProps) {
  const [openedControls, setOpenedControls] = useState(false)

  function handleClick(e: MouseEvent) {
    e.stopPropagation()
    setOpenedControls(true)
  }

  function handleCancel(e: MouseEvent) {
    e.stopPropagation()
    if (onCancel) onCancel()
    setOpenedControls(false)
  }

  function handleConfirm(e: MouseEvent) {
    e.stopPropagation()
    if (onConfirm) onConfirm()
    setOpenedControls(false)
  }

  return (
    <Menu.Item
      component={'div'}
      onClick={handleClick}
      closeMenuOnClick={false}
      style={(theme) => ({
        ...(openedControls && {
          transition: '0.325s',
          cursor: 'default',
          color: theme.colors.gray[5],
          backgroundColor: 'transparent'
        })
      })}
      {...props}
    >
      <Box pos={'relative'}>
        {children}
        <Transition
          mounted={openedControls}
          transition="fade-left"
          duration={325}
          timingFunction="ease"
        >
          {(styles) => (
            <Group
              gap={0}
              pos={'absolute'}
              pl={'md'}
              py={4}
              top={-9}
              right={-5}
              style={(theme) => ({
                ...styles,
                borderRadius: '16px',
                background: `linear-gradient(to right, transparent, ${theme.white} 27%)`
              })}
            >
              <ActionIcon
                aria-label={'cancel'}
                variant={'subtle'}
                size={26}
                radius={'50%'}
                disabled={isLoading}
                onClick={handleCancel}
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
                <IconX size={14} />
              </ActionIcon>
              <ActionIcon
                aria-label={'confirm'}
                variant={'subtle'}
                size={26}
                radius={'50%'}
                c={'green'}
                loading={isLoading}
                onClick={handleConfirm}
                sx={(theme) => ({
                  '&:hover': { backgroundColor: alpha(theme.colors.green[2], 0.35) }
                })}
              >
                <IconCheck size={14} />
              </ActionIcon>
            </Group>
          )}
        </Transition>
      </Box>
    </Menu.Item>
  )
}

export default MenuItemConfirmation
