import { ReactElement } from 'react'

interface SidebarProps {
  hidden: boolean
}

function Sidebar({ hidden }: SidebarProps): ReactElement {
  if (hidden) return <></>

  return (
    <div style={{ width: '200px' }}>
      <p>This is the navigation bar</p>
    </div>
  )
}

export default Sidebar
