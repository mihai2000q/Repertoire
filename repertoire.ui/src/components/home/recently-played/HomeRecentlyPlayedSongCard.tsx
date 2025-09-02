import Song from '../../../types/models/Song.ts'
import { useAppDispatch } from '../../../state/store.ts'
import { useNavigate } from 'react-router-dom'
import { useDisclosure } from '@mantine/hooks'
import { openArtistDrawer, openSongDrawer } from '../../../state/slice/globalSlice.ts'
import { ContextMenu } from '../../@ui/menu/ContextMenu.tsx'
import { IconDisc, IconEye, IconUser } from '@tabler/icons-react'
import OpenLinksMenuItem from '../../@ui/menu/item/song/OpenLinksMenuItem.tsx'
import YoutubeModal from '../../@ui/modal/YoutubeModal.tsx'
import { MouseEvent } from 'react'
import HomeRecentlyPlayedCard from './HomeRecentlyPlayedCard.tsx'
import CustomIconMusicNoteEighth from '../../@ui/icons/CustomIconMusicNoteEighth.tsx'

interface HomeRecentlyPlayedSongCardProps {
  song: Song
}

function HomeRecentlyPlayedSongCard({ song }: HomeRecentlyPlayedSongCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [openedMenu, { toggle: toggleMenu }] = useDisclosure(false)
  const [openedYoutube, { open: openYoutube, close: closeYoutube }] = useDisclosure(false)

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  function handleArtistClick(e: MouseEvent) {
    e.stopPropagation()
    dispatch(openArtistDrawer(song.artist.id))
  }

  function handleViewDetails() {
    navigate(`/song/${song.id}`)
  }

  function handleViewArtist() {
    navigate(`/artist/${song.artist.id}`)
  }

  function handleViewAlbum() {
    navigate(`/album/${song.album.id}`)
  }

  return (
    <ContextMenu opened={openedMenu} onChange={toggleMenu}>
      <ContextMenu.Target>
        <HomeRecentlyPlayedCard
          aria-label={`song-card-${song.title}`}
          imageUrl={song.imageUrl ?? song.album?.imageUrl}
          title={song.title}
          progress={song.progress}
          lastPlayed={song.lastTimePlayed}
          openedMenu={openedMenu}
          defaultIcon={
            <CustomIconMusicNoteEighth aria-label={`default-icon-${song.title}`} size={18} />
          }
          onClick={handleClick}
          additionalText={
            song.artist && {
              content: song.artist.name,
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
          disabled={!song.artist}
          onClick={handleViewArtist}
        >
          View Artist
        </ContextMenu.Item>
        <ContextMenu.Item
          leftSection={<IconDisc size={14} />}
          disabled={!song.album}
          onClick={handleViewAlbum}
        >
          View Album
        </ContextMenu.Item>
        <OpenLinksMenuItem song={song} openYoutube={openYoutube} />
      </ContextMenu.Dropdown>

      <YoutubeModal
        title={song.title}
        link={song.youtubeLink}
        opened={openedYoutube}
        onClose={closeYoutube}
      />
    </ContextMenu>
  )
}

export default HomeRecentlyPlayedSongCard
