import {ReactNode, useState} from 'react'
import { ActionIcon, Box, Menu, Tooltip } from '@mantine/core'
import { IconDots, IconPencil } from '@tabler/icons-react'
import { useHover } from '@mantine/hooks'

interface HeaderPanelCardProps {
  children: ReactNode
  onEditClick: () => void
  menuDropdown: ReactNode
  hideIcons?: boolean
}

function HeaderPanelCard({ children, onEditClick, menuDropdown, hideIcons }: HeaderPanelCardProps) {
  const { ref, hovered } = useHover()

  const [isMenuOpened, setIsMenuOpened] = useState(false)

  return (
    <Box ref={ref} pos={'relative'}>
      {children}

      {hideIcons !== true && (
        <Box pos={'absolute'} right={0} top={0} p={0}>
          <Menu opened={isMenuOpened} onChange={setIsMenuOpened} position={'bottom-end'} shadow={'md'}>
            <Menu.Target>
              <ActionIcon
                variant={'grey'}
                style={{ transition: '0.25s', opacity: hovered || isMenuOpened ? 1 : 0 }}
              >
                <IconDots size={18} />
              </ActionIcon>
            </Menu.Target>

            <Menu.Dropdown>{menuDropdown}</Menu.Dropdown>
          </Menu>
        </Box>
      )}

      {hideIcons !== true && (
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
      )}
    </Box>
  )
}

export default HeaderPanelCard
