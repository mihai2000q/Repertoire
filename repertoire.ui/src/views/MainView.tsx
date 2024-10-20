import { ReactElement } from 'react'
import Sidebar from '../components/Sidebar'
import Topbar from '../components/Topbar'
import { Outlet } from 'react-router-dom'
import useErrorRedirection from '../hooks/useErrorRedirection'
import { AppShell } from '@mantine/core'

function MainView(): ReactElement {
  useErrorRedirection()

  return (
    <AppShell
      layout={'alt'}
      header={{ height: 50 }}
      navbar={{
        width: 250,
        breakpoint: 'xs',
        collapsed: { mobile: true, desktop: false }
      }}
      px={'xl'}
      style={{ width: '100%' }}
    >
      <Topbar />
      <Sidebar />
      <AppShell.Main h={'100%'}>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  )
}

export default MainView
