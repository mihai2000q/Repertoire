import { ReactElement } from 'react'
import Sidebar from '../components/main/Sidebar'
import Topbar from '../components/main/Topbar'
import { Outlet } from 'react-router-dom'
import useErrorRedirection from '../hooks/useErrorRedirection'
import { AppShell, Box } from '@mantine/core'
import TitleBar from "../components/main/TitleBar.tsx";
import useAuth from "../hooks/useAuth.ts";
import useIsDesktop from 'src/hooks/useIsDesktop.ts'
import useTitleBarHeight from 'src/hooks/useTitleBarHeight.ts'

function Main(): ReactElement {
  useErrorRedirection()

  const isDesktop = useIsDesktop()
  const titleBarHeight = useTitleBarHeight()

  return (
    <Box w={'100%'} h={'100%'}>
      {isDesktop && <TitleBar />}
      <AppShell
        layout={'alt'}
        header={{ height: 50 }}
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
    </Box>
  )
}

export default Main
