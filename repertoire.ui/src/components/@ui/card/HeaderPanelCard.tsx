import { forwardRef, ReactNode, useState } from 'react'
import { ActionIcon, Box, Menu, Tooltip } from '@mantine/core'
import { IconDots, IconPencil } from '@tabler/icons-react'
import { useHover, useMergedRef } from '@mantine/hooks'

interface HeaderPanelCardProps {
  children: ReactNode
  onEditClick: () => void
  menuDropdown: ReactNode
  hideIcons?: boolean
}

const HeaderPanelCard = forwardRef<HTMLDivElement, HeaderPanelCardProps>(
  ({ children, onEditClick, menuDropdown, hideIcons }, ref) => {
    const { ref: hoverRef, hovered } = useHover()
    const mergedRef = useMergedRef(ref, hoverRef)

    const [isMenuOpened, setIsMenuOpened] = useState(false)

    return (
      <Box aria-label={'header-panel-card'} ref={mergedRef} pos={'relative'}>
        {children}

        {hideIcons !== true && (
          <Box pos={'absolute'} right={0} top={0} p={0}>
            <Menu
              opened={isMenuOpened}
              onChange={setIsMenuOpened}
              position={'bottom-end'}
              shadow={'md'}
            >
              <Menu.Target>
                <ActionIcon
                  aria-label={'more-menu'}
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
                aria-label={'edit-header'}
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
)

HeaderPanelCard.displayName = 'HeaderPanelCard'

export default HeaderPanelCard
