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
import SongDrawer from '../components/main/SongDrawer.tsx'
import { useAppDispatch, useAppSelector } from '../state/store.ts'
import { closeSongDrawer as closeSongDrawerRedux, closeAlbumDrawer as closeAlbumDrawerRedux } from '../state/globalSlice.ts'
import AlbumDrawer from "../components/main/AlbumDrawer.tsx";

function Main(): ReactElement {
  useErrorRedirection()

  const isDesktop = useIsDesktop()
  const titleBarHeight = useTitleBarHeight()
  const dispatch = useAppDispatch()

  const openedAlbumDrawer = useAppSelector((state) => state.global.albumDrawer.open)
  const closeAlbumDrawer = () => dispatch(closeAlbumDrawerRedux())

  const openedSongDrawer = useAppSelector((state) => state.global.songDrawer.open)
  const closeSongDrawer = () => dispatch(closeSongDrawerRedux())

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

      <AlbumDrawer opened={openedAlbumDrawer} close={closeAlbumDrawer} />
      <SongDrawer opened={openedSongDrawer} close={closeSongDrawer} />
    </Box>
  )
}

export default Main
