import { ReactElement } from 'react'
import Sidebar from '../components/Sidebar'
import Topbar from '../components/Topbar'
import { Outlet } from 'react-router-dom'
import useErrorRedirection from '../hooks/useErrorRedirection'

function MainView(): ReactElement {
  useErrorRedirection()

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'row',
        width: '100%',
        height: '100%'
      }}
    >
      <Sidebar />
      <div style={{ width: '100%' }}>
        <Topbar />
        <div>
          <Outlet />
        </div>
      </div>
    </div>
  )
}

export default MainView
