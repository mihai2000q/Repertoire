import { ReactElement } from 'react'

interface TopbarProps {
  hidden: boolean
}

function Topbar({ hidden }: TopbarProps): ReactElement {
  if (hidden) return <></>

  return (
    <div>
      <p>This is the top bar</p>
    </div>
  )
}

export default Topbar
