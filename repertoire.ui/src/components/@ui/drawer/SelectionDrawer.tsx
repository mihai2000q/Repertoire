import { ActionIcon, CloseButton, Drawer, Group, Menu, Text } from '@mantine/core'
import { ReactNode } from 'react'
import { IconDots } from '@tabler/icons-react'

interface SelectionDrawerProps {
  opened: boolean
  text: string
  actionIcons: ReactNode
  onClose: () => void
  menu?: {
    opened: boolean
    toggle: () => void
    dropdown: ReactNode
  }
}

function SelectionDrawer({ opened, text, actionIcons, onClose, menu, ...props }: SelectionDrawerProps) {
  return (
    <Drawer
      opened={opened}
      onClose={() => {}}
      position={'bottom'}
      withOverlay={false}
      withCloseButton={false}
      styles={{
        inner: {
          left: '50%',
          marginBottom: '32px',
          width: 'fit-content',
          zIndex: 100
        },
        content: {
          borderRadius: '24px',
          height: 'fit-content'
        },
        body: { padding: 0 }
      }}
      {...props}
    >
      <Group gap={'xxs'} py={'xs'} px={'lg'}>
        <Text fw={500} mr={'xs'} c={'gray.6'}>
          {text}
        </Text>

        {actionIcons}

        {menu && (
          <Menu
            position={'top-start'}
            transitionProps={{ duration: 160, transition: 'pop-bottom-left' }}
            opened={menu.opened}
            onChange={menu.toggle}
          >
            <Menu.Target>
              <ActionIcon aria-label={'more-menu'} variant={'grey'}>
                <IconDots size={18} />
              </ActionIcon>
            </Menu.Target>

            {menu.dropdown}
          </Menu>
        )}

        <CloseButton aria-label={'close-drawer'} onClick={onClose} />
      </Group>
    </Drawer>
  )
}

export default SelectionDrawer
