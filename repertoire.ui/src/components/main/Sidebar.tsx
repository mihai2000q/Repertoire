import { ReactElement } from 'react'
import { AppShell, NavLink } from '@mantine/core'
import { sidebarLinks } from '../../data/sidebarLinks'
import wallpaper from '../../assets/wallpapers/sidebar.jpg'
import { useLocation, useNavigate } from 'react-router-dom'
import { createStyles } from '@mantine/emotion'

const useStyles = createStyles(() => ({
  backdrop: {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    width: '100%',
    height: '100%',
    opacity: 0.3,
    backgroundImage: `url(${wallpaper})`,
    backgroundSize: 'cover',
    backgroundPosition: '20%',

    '&::before': {
      content: '""',
      position: 'absolute',
      top: 0,
      left: 0,
      width: '100%',
      height: '100%',
      backdropFilter: 'blur(30px)'
    }
  }
}))

function Sidebar(): ReactElement {
  const location = useLocation()
  const navigate = useNavigate()

  const { classes } = useStyles()

  return (
    <AppShell.Navbar py={'xl'} px={'lg'} top={'unset'} bg={'transparent'} withBorder={false}>
      <div className={classes.backdrop} />
      <div style={{ zIndex: 2 }}>
        {sidebarLinks.map((sidebarLink) => (
          <NavLink
            key={sidebarLink.label}
            label={sidebarLink.label}
            leftSection={sidebarLink.icon}
            active={
              location.pathname === sidebarLink.link ||
              sidebarLink.subLinks.some((link) => location.pathname.startsWith(link))
            }
            onClick={() => navigate(sidebarLink.link)}
          />
        ))}
      </div>
    </AppShell.Navbar>
  )
}

export default Sidebar
