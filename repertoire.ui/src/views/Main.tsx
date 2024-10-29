import { ReactElement } from 'react'
import Sidebar from '../components/main/Sidebar'
import Topbar from '../components/main/Topbar'
import { Outlet } from 'react-router-dom'
import useErrorRedirection from '../hooks/useErrorRedirection'
import { AppShell } from '@mantine/core'
import TitleBar from "../components/main/TitleBar.tsx";
import useAuth from "../hooks/useAuth.ts";

function Main(): ReactElement {
  useErrorRedirection()

  return (
    <div style={{  }}>
      {import.meta.env.VITE_PLATFORM === 'desktop' && <TitleBar />}
      <AppShell
        layout={'alt'}
        header={{ height: 50 }}
        navbar={{
          width: 250,
          breakpoint: 'xs',
          collapsed: { mobile: true, desktop: false }
        }}
        px={'xl'}
        mt={45}
      >
        <Topbar />
        <Sidebar />
        <AppShell.Main style={{ minHeight: 0 }}>
          <Outlet />
        </AppShell.Main>
      </AppShell>
    </div>
  )
}

export default Main
