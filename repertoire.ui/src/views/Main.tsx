import { ReactElement } from 'react'
import Sidebar from '../components/main/Sidebar'
import Topbar from '../components/main/Topbar'
import { Outlet } from 'react-router-dom'
import useErrorRedirection from '../hooks/useErrorRedirection'
import { AppShell, Box, ScrollArea } from '@mantine/core'
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

function Main(): ReactElement {
  useErrorRedirection()
  useNetworkDisconnected()

  const isDesktop = useIsDesktop()
  const titleBarHeight = useTitleBarHeight()
  const topbarHeight = useTopbarHeight()

  const [mobileSidebarOpened, { toggle: toggleSidebarMobile }] = useDisclosure()

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
          <ScrollArea.Autosize
            scrollbars={'y'}
            scrollbarSize={10}
            h={`calc(100vh - ${titleBarHeight} - ${topbarHeight})`}
            styles={{
              root: {
                height: '100%'
              },
              viewport: {
                '> div': {
                  display: 'block !important', // Critical for truncation
                  minWidth: '100%', // Prevents content shrinkage
                  width: 0 // Magic that enables proper truncation
                }
              }
            }}
          >
            <Outlet />
          </ScrollArea.Autosize>
        </AppShell.Main>
      </AppShell>

      <ArtistDrawer />
      <AlbumDrawer />
      <SongDrawer />
    </Box>
  )
}

export default Main
