import Album from '../../../types/models/Album.ts'
import { useAppDispatch } from '../../../state/store.ts'
import { useNavigate } from 'react-router-dom'
import { useDisclosure } from '@mantine/hooks'
import { openAlbumDrawer, openArtistDrawer } from '../../../state/slice/globalSlice.ts'
import { ContextMenu } from '../../@ui/menu/ContextMenu.tsx'
import { IconEye, IconUser } from '@tabler/icons-react'
import { MouseEvent } from 'react'
import HomeRecentlyPlayedCard from './HomeRecentlyPlayedCard.tsx'
import CustomIconAlbumVinyl from '../../@ui/icons/CustomIconAlbumVinyl.tsx'

interface HomeRecentlyPlayedAlbumCardProps {
  album: Album
}

function HomeRecentlyPlayedAlbumCard({ album }: HomeRecentlyPlayedAlbumCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [openedMenu, { toggle: toggleMenu }] = useDisclosure(false)

  function handleClick() {
    dispatch(openAlbumDrawer(album.id))
  }

  function handleArtistClick(e: MouseEvent) {
    e.stopPropagation()
    dispatch(openArtistDrawer(album.artist.id))
  }

  function handleViewDetails() {
    navigate(`/album/${album.id}`)
  }

  function handleViewArtist() {
    navigate(`/artist/${album.artist.id}`)
  }

  return (
    <ContextMenu opened={openedMenu} onChange={toggleMenu}>
      <ContextMenu.Target>
        <HomeRecentlyPlayedCard
          aria-label={`album-card-${album.title}`}
          imageUrl={album.imageUrl}
          title={album.title}
          progress={album.progress}
          lastPlayed={album.lastTimePlayed}
          openedMenu={openedMenu}
          defaultIcon={
            <CustomIconAlbumVinyl aria-label={`default-icon-${album.title}`} size={18} />
          }
          onClick={handleClick}
          additionalText={
            album.artist && {
              content: album.artist.name,
              onClick: handleArtistClick
            }
          }
        />
      </ContextMenu.Target>

      <ContextMenu.Dropdown>
        <ContextMenu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
          View Details
        </ContextMenu.Item>
        <ContextMenu.Item
          leftSection={<IconUser size={14} />}
          disabled={!album.artist}
          onClick={handleViewArtist}
        >
          View Artist
        </ContextMenu.Item>
      </ContextMenu.Dropdown>
    </ContextMenu>
  )
}

export default HomeRecentlyPlayedAlbumCard
