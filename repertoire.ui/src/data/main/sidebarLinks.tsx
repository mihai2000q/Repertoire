import { ReactElement } from 'react'
import {
  IconAlbum,
  IconHomeFilled,
  IconMusic,
  IconPlaylist,
  IconUserFilled
} from '@tabler/icons-react'

interface SidebarLink {
  icon: ReactElement
  label: string
  link: string,
  subLinks: string[]
}

export const sidebarLinks: SidebarLink[] = [
  { icon: <IconHomeFilled />, label: 'Home', link: '/home', subLinks: [] },
  { icon: <IconUserFilled />, label: 'Artists', link: '/artists', subLinks: ['/artist'] },
  { icon: <IconAlbum />, label: 'Albums', link: '/albums', subLinks: ['/album'] },
  { icon: <IconMusic stroke={1.75} />, label: '/songs', link: '/songs', subLinks: ['/song'] },
  { icon: <IconPlaylist stroke={1.75} />, label: '/playlists', link: '/playlists', subLinks: ['/playlist'] }
]
