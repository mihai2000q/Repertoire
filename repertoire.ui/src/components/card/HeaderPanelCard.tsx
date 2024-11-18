import { ReactNode } from 'react'
import {ActionIcon, Box, Menu, Tooltip} from '@mantine/core'
import {IconDots, IconPencil} from '@tabler/icons-react'
import { useHover } from '@mantine/hooks'

interface HeaderPanelCardProps {
  children: ReactNode
  onEditClick: () => void,
  menuDropdown: ReactNode
}

function HeaderPanelCard({ children, onEditClick, menuDropdown }: HeaderPanelCardProps) {
  const { ref, hovered } = useHover()

  return (
    <Box ref={ref} pos={'relative'}>
      {children}

      <Box pos={'absolute'} right={0} top={0} p={0}>
        <Menu position={'bottom-end'} shadow={'md'}>
          <Menu.Target>
            <ActionIcon
              variant={'grey'}
              style={{ transition: '0.25s', opacity: hovered ? 1 : 0 }}
            >
              <IconDots size={18} />
            </ActionIcon>
          </Menu.Target>

          <Menu.Dropdown>
            {menuDropdown}
          </Menu.Dropdown>
        </Menu>
      </Box>

      <Box pos={'absolute'} right={0} bottom={-12} p={0}>
        <Tooltip label={'Edit Header'} openDelay={500}>
          <ActionIcon
            variant={'grey'}
            style={{ transition: '0.25s', opacity: hovered ? 1 : 0 }}
            onClick={onEditClick}
          >
            <IconPencil size={18} />
          </ActionIcon>
        </Tooltip>
      </Box>
    </Box>
  )
}

export default HeaderPanelCard