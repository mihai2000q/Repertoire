import { ReactElement } from 'react'
import { AppShell, NavLink } from '@mantine/core'
import { sidebarLinks } from '../../data/sidebarLinks'
import { useLocation, useNavigate } from 'react-router-dom'

function Sidebar(): ReactElement {
  const location = useLocation()
  const navigate = useNavigate()

  return (
    <AppShell.Navbar py={'xl'} px={'lg'}>
      {sidebarLinks.map((sidebarLink) => (
        <NavLink
          key={sidebarLink.label}
          label={sidebarLink.label}
          leftSection={sidebarLink.icon}
          active={location.pathname === sidebarLink.link}
          onClick={() => navigate(sidebarLink.link)}
        />
      ))}
    </AppShell.Navbar>
  )
}

export default Sidebar
