import { ReactElement } from 'react'
import { AppShell, NavLink } from '@mantine/core'
import { sidebarLinks } from '../data/sidebarLinks'
import { useLocation } from 'react-router-dom'

function Sidebar(): ReactElement {
  const location = useLocation()

  return (
    <AppShell.Navbar py={'xl'} px={'lg'}>
      {sidebarLinks.map((sidebarLink) => (
        <NavLink
          key={sidebarLink.label}
          label={sidebarLink.label}
          leftSection={sidebarLink.icon}
          active={location.pathname === sidebarLink.link}
        />
      ))}
    </AppShell.Navbar>
  )
}

export default Sidebar
