import { ReactElement } from 'react'
import { ActionIcon, AppShell, Group, Space } from '@mantine/core'
import { IconBellFilled, IconChevronLeft, IconChevronRight, IconMenu2 } from '@tabler/icons-react'
import { useWindowScroll } from '@mantine/hooks'
import { useNavigate } from 'react-router-dom'
import useIsDesktop from '../../hooks/useIsDesktop.ts'
import TopbarSearch from './TopbarSearch.tsx'
import TopbarUser from './TopbarUser.tsx'

interface TopbarProps {
  toggleSidebar: () => void
}

function Topbar({ toggleSidebar }: TopbarProps): ReactElement {
  const navigate = useNavigate()
  const isDesktop = useIsDesktop()
  const [scrollPosition] = useWindowScroll()

  function handleGoBack() {
    navigate(-1)
  }

  function handleGoForward() {
    navigate(1)
  }

  return (
    <AppShell.Header
      px={'md'}
      withBorder={false}
      top={'unset'}
      style={(theme) => ({
        transition: '0.35s',
        ...(scrollPosition.y !== 0 && { boxShadow: theme.shadows.md })
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
          sx={(theme) => ({
            order: 1,
            [`@media(max-width: ${theme.breakpoints.sm})`]: {
              order: 3,
              justifySelf: 'center'
            }
          })}
        />

        <Space hiddenFrom={'sm'} flex={1} style={{ order: 2 }} />

        {isDesktop && (
          <Group
            gap={0}
            ml={'xs'}
            sx={(theme) => ({
              order: 3,
              [`@media(max-width: ${theme.breakpoints.sm})`]: { order: 1 }
            })}
          >
            <ActionIcon
              aria-label={'back-button'}
              size={'lg'}
              variant={'grey'}
              radius={'50%'}
              disabled={window.history.state?.idx < 1}
              onClick={handleGoBack}
            >
              <IconChevronLeft size={20} />
            </ActionIcon>

            <ActionIcon
              aria-label={'forward-button'}
              size={'lg'}
              variant={'grey'}
              radius={'50%'}
              disabled={window.history.state?.idx >= window.history.length - 1}
              onClick={handleGoForward}
            >
              <IconChevronRight size={20} />
            </ActionIcon>
          </Group>
        )}

        <Space flex={1} style={{ order: 4 }} />

        <ActionIcon
          variant={'subtle'}
          size={'lg'}
          sx={(theme) => ({
            borderRadius: '50%',
            color: theme.colors.gray[6],
            '&:hover': {
              boxShadow: theme.shadows.sm,
              backgroundColor: theme.colors.primary[0],
              color: theme.colors.primary[6]
            }
          })}
          style={{ order: 5 }}
        >
          <IconBellFilled size={18} />
        </ActionIcon>

        <TopbarUser style={{ order: 6 }} />
      </Group>
    </AppShell.Header>
  )
}

export default Topbar
