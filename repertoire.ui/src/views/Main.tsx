import { ReactElement } from 'react'
import Sidebar from '../components/main/Sidebar'
import Topbar from '../components/main/Topbar'
import { Outlet } from 'react-router-dom'
import useErrorRedirection from '../hooks/useErrorRedirection'
import { AppShell, Box } from '@mantine/core'
import TitleBar from '../components/main/TitleBar'
import useAuth from '../hooks/useAuth'
import useIsDesktop from '../hooks/useIsDesktop'
import useTitleBarHeight from '../hooks/useTitleBarHeight'
import SongDrawer from '../components/main/drawer/SongDrawer.tsx'
import AlbumDrawer from '../components/main/drawer/AlbumDrawer.tsx'
import ArtistDrawer from '../components/main/drawer/ArtistDrawer.tsx'

function Main(): ReactElement {
  useErrorRedirection()

  const isDesktop = useIsDesktop()
  const titleBarHeight = useTitleBarHeight()

  return (
    <Box w={'100%'} h={'100%'}>
      {isDesktop && <TitleBar />}
      <AppShell
        layout={'alt'}
        header={{ height: 65 }}
        navbar={{
          width: 250,
          breakpoint: 'xs',
          collapsed: { mobile: true, desktop: false }
        }}
        px={'xl'}
        w={'100%'}
        h={`calc(100% - ${titleBarHeight})`}
        mt={titleBarHeight}
        disabled={!useAuth()}
      >
        <Topbar />
        <Sidebar />
        <AppShell.Main h={'100%'} mih={0}>
          <Outlet />
        </AppShell.Main>
      </AppShell>

      <ArtistDrawer />
      <AlbumDrawer />
      <SongDrawer />
    </Box>
  )
}

export default Main
