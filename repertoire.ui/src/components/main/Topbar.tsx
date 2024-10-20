import { ReactElement } from 'react'
import { AppShell } from '@mantine/core'

function Topbar(): ReactElement {
  return (
    <AppShell.Header px={'md'} withBorder={false}>
      <p>This is the top bar</p>
    </AppShell.Header>
  )
}

export default Topbar
