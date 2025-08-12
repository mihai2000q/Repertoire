import Artist from '../../../types/models/Artist.ts'
import { useAppDispatch } from '../../../state/store.ts'
import { useNavigate } from 'react-router-dom'
import { useDisclosure } from '@mantine/hooks'
import { openArtistDrawer } from '../../../state/slice/globalSlice.ts'
import { ContextMenu } from '../../@ui/menu/ContextMenu.tsx'
import { IconEye, IconUser } from '@tabler/icons-react'
import HomeRecentlyPlayedCard from './HomeRecentlyPlayedCard.tsx'

interface HomeRecentlyPlayedArtistCardProps {
  artist: Artist
}

function HomeRecentlyPlayedArtistCard({ artist }: HomeRecentlyPlayedArtistCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [openedMenu, { toggle: toggleMenu }] = useDisclosure(false)

  function handleClick() {
    dispatch(openArtistDrawer(artist.id))
  }

  function handleViewDetails() {
    navigate(`/artist/${artist.id}`)
  }

  return (
    <ContextMenu shadow={'lg'} opened={openedMenu} onChange={toggleMenu}>
      <ContextMenu.Target>
        <HomeRecentlyPlayedCard
          aria-label={`artist-card-${artist.name}`}
          imageUrl={artist.imageUrl}
          title={artist.name}
          progress={artist.progress}
          lastPlayed={artist.lastTimePlayed}
          openedMenu={openedMenu}
          defaultIcon={<IconUser aria-label={`default-icon-${artist.name}`} size={18} />}
          onClick={handleClick}
          isArtist={true}
        />
      </ContextMenu.Target>

      <ContextMenu.Dropdown>
        <ContextMenu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
          View Details
        </ContextMenu.Item>
      </ContextMenu.Dropdown>
    </ContextMenu>
  )
}

export default HomeRecentlyPlayedArtistCard
