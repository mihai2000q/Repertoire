import Song from '../../../types/models/Song.ts'
import { Avatar, Center, Stack, Text } from '@mantine/core'
import { useState } from 'react'
import { openArtistDrawer, openSongDrawer } from '../../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../../state/store.ts'
import CustomIconMusicNote from '../../@ui/icons/CustomIconMusicNote.tsx'
import { IconDisc, IconEye, IconUser } from '@tabler/icons-react'
import OpenLinksMenuItem from '../../@ui/menu/item/song/OpenLinksMenuItem.tsx'
import { useDisclosure } from '@mantine/hooks'
import YoutubeModal from '../../@ui/modal/YoutubeModal.tsx'
import { useNavigate } from 'react-router-dom'
import { ContextMenu } from '../../@ui/menu/ContextMenu.tsx'

interface HomeSongCardProps {
  song: Song
}

function HomeSongCard({ song }: HomeSongCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [isImageHovered, setIsImageHovered] = useState(false)
  const [openedMenu, { toggle: toggleMenu }] = useDisclosure(false)
  const [openedYoutube, { open: openYoutube, close: closeYoutube }] = useDisclosure(false)

  const isSelected = isImageHovered || openedMenu

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  function handleArtistClick() {
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
    <Stack
      aria-label={`song-card-${song.title}`}
      align={'center'}
      gap={0}
      style={{
        transition: '0.25s',
        ...(isSelected && { transform: 'scale(1.05)' })
      }}
      w={'max(10vw, 150px)'}
    >
      <ContextMenu shadow={'lg'} opened={openedMenu} onChange={toggleMenu}>
        <ContextMenu.Target>
          <Avatar
            radius={'10%'}
            w={'100%'}
            h={'unset'}
            src={song.imageUrl ?? song.album?.imageUrl}
            alt={(song.imageUrl ?? song.album?.imageUrl) && song.title}
            bg={'gray.5'}
            sx={(theme) => ({
              aspectRatio: 1,
              cursor: 'pointer',
              transition: '0.25s',
              boxShadow: theme.shadows.xl,
              ...(isSelected && { boxShadow: theme.shadows.xxl })
            })}
            onMouseEnter={() => setIsImageHovered(true)}
            onMouseLeave={() => setIsImageHovered(false)}
            onClick={handleClick}
          >
            <Center c={'white'}>
              <CustomIconMusicNote aria-label={`default-icon-${song.title}`} size={50} />
            </Center>
          </Avatar>
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
      </ContextMenu>

      <Stack w={'100%'} pt={'xs'} gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={600} lineClamp={2} ta={'center'}>
          {song.title}
        </Text>
        {song.artist && (
          <Text
            fw={500}
            ta={'center'}
            c={'dimmed'}
            truncate={'end'}
            onClick={handleArtistClick}
            sx={{
              cursor: 'pointer',
              '&:hover': { textDecoration: 'underline' }
            }}
          >
            {song.artist.name}
          </Text>
        )}
      </Stack>

      <YoutubeModal
        title={song.title}
        link={song.youtubeLink}
        opened={openedYoutube}
        onClose={closeYoutube}
      />
    </Stack>
  )
}

export default HomeSongCard
