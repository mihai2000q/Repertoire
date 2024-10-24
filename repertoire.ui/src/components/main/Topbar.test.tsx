import { mantineRender } from '../../test-utils.tsx'
import Topbar from './Topbar.tsx'
import { AppShell } from '@mantine/core'

describe('Topbar', () => {
  const render = () => mantineRender(
    <AppShell>
      <Topbar />
    </AppShell>
  )

  it('should render', () => {
    render()
  })
})
