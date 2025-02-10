import { ReactElement } from 'react'
import { IconHomeFilled, IconUserFilled } from '@tabler/icons-react'
import CustomIconMusicNote from '../../components/@ui/icons/CustomIconMusicNote.tsx'
import CustomIconAlbumVinyl from '../../components/@ui/icons/CustomIconAlbumVinyl.tsx'
import CustomIconPlaylist2 from '../../components/@ui/icons/CustomIconPlaylist2.tsx'

interface SidebarLink {
  icon: ReactElement
  label: string
  link: string
  subLinks: string[]
}

export const sidebarLinks: SidebarLink[] = [
  { icon: <IconHomeFilled />, label: 'Home', link: '/home', subLinks: [] },
  { icon: <IconUserFilled />, label: 'Artists', link: '/artists', subLinks: ['/artist'] },
  { icon: <CustomIconAlbumVinyl />, label: 'Albums', link: '/albums', subLinks: ['/album'] },
  { icon: <CustomIconMusicNote />, label: 'Songs', link: '/songs', subLinks: ['/song'] },
  {
    icon: <CustomIconPlaylist2 />,
    label: 'Playlists',
    link: '/playlists',
    subLinks: ['/playlist']
  }
]
