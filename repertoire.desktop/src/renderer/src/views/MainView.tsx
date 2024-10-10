import { ReactElement } from 'react'
import Sidebar from '@renderer/components/Sidebar'
import Topbar from '@renderer/components/Topbar'
import useAuth from '@renderer/hooks/useAuth'
import { Outlet } from 'react-router-dom'

function MainView(): ReactElement {
  const isLayoutHidden = !useAuth()

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'row',
        width: '100%',
        height: '100%'
      }}
    >
      <Sidebar hidden={isLayoutHidden} />
      <div style={{ width: '100%' }}>
        <Topbar hidden={isLayoutHidden} />
        <div>
          <Outlet />
        </div>
      </div>
    </div>
  )
}

export default MainView
