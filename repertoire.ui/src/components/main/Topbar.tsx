import { ReactElement } from 'react'
import { ActionIcon, AppShell, Flex, Group, Space, useMantineTheme } from '@mantine/core'
import { IconBellFilled, IconMenu2 } from '@tabler/icons-react'
import { useMediaQuery } from '@mantine/hooks'
import useIsDesktop from '../../hooks/useIsDesktop.ts'
import TopbarSearch from './topbar/TopbarSearch.tsx'
import TopbarUser from './topbar/TopbarUser.tsx'
import useMainScroll from '../../hooks/useMainScroll.ts'
import TopbarNavigation from './topbar/TopbarNavigation.tsx'

interface TopbarProps {
  toggleSidebar: () => void
}

function Topbar({ toggleSidebar }: TopbarProps): ReactElement {
  const isDesktop = useIsDesktop()
  const { isTopScrollPositionOver0 } = useMainScroll()
  const shiftOrder = isDesktop ? 0 : 1

  const theme = useMantineTheme()
  const isSmallScreen = useMediaQuery(`(max-width: ${theme.breakpoints.sm})`)

  return (
    <AppShell.Header
      px={'md'}
      withBorder={false}
      top={'unset'}
      style={(theme) => ({
        transition: '0.35s',
        ...(isTopScrollPositionOver0 && { boxShadow: theme.shadows.md })
      })}
    >
      <Group h={'100%'} gap={0}>
        <ActionIcon
          hiddenFrom={'sm'}
          aria-label={'toggle-sidebar'}
          variant={'grey'}
          size={'lg'}
          onClick={toggleSidebar}
          style={{ order: 0 }}
        >
          <IconMenu2 />
        </ActionIcon>

        <TopbarSearch
          w={'max(15vw, 200px)'}
          comboboxProps={{
            width: 'max(20vw, 325px)',
            position: isSmallScreen ? 'bottom' : 'bottom-start'
          }}
          dropdownMinHeight={'max(26vh, 200px)'}
          style={{ order: isSmallScreen ? 3 - shiftOrder : 1 }}
        />

        <Space hiddenFrom={'sm'} flex={1} style={{ order: 2 }} />

        {isDesktop && (
          <Flex style={{ order: isSmallScreen ? 1 : 3 }}>
            <TopbarNavigation />
          </Flex>
        )}

        <Space flex={1} style={{ order: 4 - shiftOrder }} />

        <ActionIcon
          variant={'grey-primary'}
          size={'lg'}
          radius={'50%'}
          sx={(theme) => ({
            '&:hover': { boxShadow: theme.shadows.sm }
          })}
          style={{ order: 5 - shiftOrder }}
        >
          <IconBellFilled size={18} />
        </ActionIcon>

        <TopbarUser style={{ order: 6 - shiftOrder }} />
      </Group>
    </AppShell.Header>
  )
}

export default Topbar
