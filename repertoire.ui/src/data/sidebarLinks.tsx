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
  link: string
}

export const sidebarLinks: SidebarLink[] = [
  { icon: <IconHomeFilled />, label: 'Home', link: '/home' },
  { icon: <IconUserFilled />, label: 'Artists', link: '/artists' },
  { icon: <IconAlbum />, label: 'Albums', link: '/albums' },
  { icon: <IconMusic stroke={1.75} />, label: 'Songs', link: '/songs' },
  { icon: <IconPlaylist stroke={1.75} />, label: 'Playlists', link: '/playlists' }
]
