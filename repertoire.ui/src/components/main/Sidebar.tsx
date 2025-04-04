import { ReactElement } from 'react'
import {
  ActionIcon,
  AppShell,
  Box,
  Group,
  NavLink,
  Stack,
  Title,
  useMantineTheme
} from '@mantine/core'
import { sidebarLinks } from '../../data/main/sidebarLinks.tsx'
import wallpaper from '../../assets/wallpapers/sidebar.jpg'
import { useLocation, useNavigate } from 'react-router-dom'
import { createStyles } from '@mantine/emotion'
import { IconLayoutSidebarLeftCollapseFilled } from '@tabler/icons-react'
import { useMediaQuery } from '@mantine/hooks'

const useStyles = createStyles(() => ({
  backdrop: {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    width: '100%',
    height: '100%',
    opacity: 0.3,
    backgroundImage: `url(${wallpaper})`,
    backgroundSize: 'cover',
    backgroundPosition: '20%',

    '&::before': {
      content: '""',
      position: 'absolute',
      top: 0,
      left: 0,
      width: '100%',
      height: '100%',
      backdropFilter: 'blur(30px)'
    }
  }
}))

interface SidebarProps {
  toggleSidebarOnMobile: () => void
}

function Sidebar({ toggleSidebarOnMobile }: SidebarProps): ReactElement {
  const location = useLocation()
  const navigate = useNavigate()

  const theme = useMantineTheme()
  const isMobile = useMediaQuery(`(max-width: ${theme.breakpoints.sm})`)

  const { classes } = useStyles()

  return (
    <AppShell.Navbar
      py={'xl'}
      px={'lg'}
      top={'unset'}
      bg={{ base: 'white', sm: 'transparent' }}
      withBorder={false}
    >
      <Box visibleFrom={'sm'}>
        <div className={classes.backdrop} />
      </Box>
      <Stack pos={'relative'} gap={0}>
        <Group hiddenFrom={'sm'} pos={'relative'} mb={'lg'}>
          <ActionIcon
            aria-label={'toggle-sidebar'}
            variant={'grey'}
            size={'xl'}
            pos={'absolute'}
            top={-7.5}
            left={20}
            onClick={toggleSidebarOnMobile}
          >
            <IconLayoutSidebarLeftCollapseFilled />
          </ActionIcon>

          <Title w={'100%'} order={6} ta={'center'} c={'dimmed'} fw={800}>
            Navigation
          </Title>
        </Group>

        <Stack align={'center'} gap={2}>
          {sidebarLinks.map((sidebarLink) => (
            <NavLink
              key={sidebarLink.label}
              role={'link'}
              label={sidebarLink.label}
              leftSection={sidebarLink.icon}
              active={
                location.pathname === sidebarLink.link ||
                sidebarLink.subLinks.some((link) => location.pathname.startsWith(link))
              }
              onClick={() => {
                navigate(sidebarLink.link)
                if (isMobile) toggleSidebarOnMobile()
              }}
            />
          ))}
        </Stack>
      </Stack>
    </AppShell.Navbar>
  )
}

export default Sidebar
