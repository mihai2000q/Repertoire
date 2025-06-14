import { ReactElement } from 'react'
import Sidebar from '../components/main/Sidebar'
import Topbar from '../components/main/Topbar'
import { Outlet } from 'react-router-dom'
import useErrorRedirection from '../hooks/useErrorRedirection'
import { alpha, AppShell, Box } from '@mantine/core'
import TitleBar from '../components/main/TitleBar'
import useAuth from '../hooks/useAuth'
import useIsDesktop from '../hooks/useIsDesktop'
import useTitleBarHeight from '../hooks/useTitleBarHeight'
import SongDrawer from '../components/main/drawer/SongDrawer.tsx'
import AlbumDrawer from '../components/main/drawer/AlbumDrawer.tsx'
import ArtistDrawer from '../components/main/drawer/ArtistDrawer.tsx'
import { useDisclosure } from '@mantine/hooks'
import useNetworkDisconnected from '../hooks/useNetworkDisconnected.tsx'
import useTopbarHeight from '../hooks/useTopbarHeight.ts'
import useMainScroll from '../hooks/useMainScroll.ts'
import { createStyles } from '@mantine/emotion'
import PlaylistDrawer from '../components/main/drawer/PlaylistDrawer.tsx'

const useStyles = createStyles((theme) => ({
  scrollbar: {
    '&::-webkit-scrollbar': {
      width: 9
    },

    '&::-webkit-scrollbar-track-piece': {
      backgroundColor: 'transparent',
      '&:hover': {
        backgroundColor: theme.colors.gray[2]
      }
    },

    '&::-webkit-scrollbar-thumb': {
      borderRadius: theme.radius.md,
      backgroundColor: alpha(theme.colors.gray[6], 0.9),

      '&:hover': {
        backgroundColor: alpha(theme.colors.gray[7], 0.75)
      }
    }
  }
}))

function Main(): ReactElement {
  useErrorRedirection()
  useNetworkDisconnected()

  const isDesktop = useIsDesktop()
  const titleBarHeight = useTitleBarHeight()
  const topbarHeight = useTopbarHeight()

  const [mobileSidebarOpened, { toggle: toggleSidebarMobile }] = useDisclosure()

  const { ref: viewportRef } = useMainScroll()
  const { classes } = useStyles()

  return (
    <Box w={'100%'} h={'100%'}>
      {isDesktop && <TitleBar />}
      <AppShell
        layout={'alt'}
        header={{ height: topbarHeight }}
        navbar={{
          width: 'max(15vw, 250px)',
          breakpoint: 'sm',
          collapsed: { mobile: !mobileSidebarOpened, desktop: false }
        }}
        w={'100%'}
        h={`calc(100% - ${titleBarHeight})`}
        mt={titleBarHeight}
        disabled={!useAuth()}
      >
        <Topbar toggleSidebar={toggleSidebarMobile} />
        <Sidebar toggleSidebarOnMobile={toggleSidebarMobile} />
        <AppShell.Main h={'100%'} mih={0}>
          <Box
            ref={viewportRef}
            className={classes.scrollbar}
            style={{ height: '100%', overflow: 'auto' }}
          >
            <Outlet />
          </Box>
        </AppShell.Main>
      </AppShell>

      <ArtistDrawer />
      <AlbumDrawer />
      <SongDrawer />
      <PlaylistDrawer />
    </Box>
  )
}

export default Main
