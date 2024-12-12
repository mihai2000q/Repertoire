import Song from '../../types/models/Song.ts'
import { ActionIcon, alpha, Avatar, Group, Menu, Stack, Text } from '@mantine/core'
import songPlaceholder from '../../assets/image-placeholder-1.jpg'
import { useAppDispatch } from '../../state/store.ts'
import { openSongDrawer } from '../../state/globalSlice.ts'
import { useHover } from '@mantine/hooks'
import { MouseEvent, useState } from 'react'
import { IconDots, IconTrash } from '@tabler/icons-react'

interface PlaylistSongCardProps {
  song: Song
  handleRemove: () => void
}

function PlaylistSongCard({ song, handleRemove }: PlaylistSongCardProps) {
  const dispatch = useAppDispatch()
  const { ref, hovered } = useHover()

  const [isMenuOpened, setIsMenuOpened] = useState(false)

  const isSelected = hovered || isMenuOpened

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  function handleDelete(e: MouseEvent) {
    e.stopPropagation()
    handleRemove()
  }

  return (
    <Group
      ref={ref}
      align={'center'}
      wrap={'nowrap'}
      sx={(theme) => ({
        cursor: 'default',
        transition: '0.3s',
        '&:hover': {
          boxShadow: theme.shadows.xl,
          backgroundColor: alpha(theme.colors.cyan[0], 0.15)
        }
      })}
      px={'md'}
      py={'xs'}
      onClick={handleClick}
    >
      <Text fw={500} w={35} ta={'center'}>
        {song.playlistTrackNo}
      </Text>

      <Avatar radius={'8px'} src={song.imageUrl ?? songPlaceholder} />

      <Stack flex={1} gap={0} style={{ overflow: 'hidden' }}>
        <Group gap={4}>
          <Text fw={500} truncate={'end'}>{song.title}</Text>
          {song.album && <Text fz={'sm'} c={'dimmed'} truncate={'end'}>- {song.album.title}</Text>}
        </Group>
        {song.artist && <Text fz={'sm'} c={'dimmed'}>{song.artist.name}</Text>}
      </Stack>

      <Menu position={'bottom-end'} opened={isMenuOpened} onChange={setIsMenuOpened}>
        <Menu.Target>
          <ActionIcon
            size={'md'}
            variant={'grey'}
            onClick={(e) => e.stopPropagation()}
            style={{
              transition: '0.3s',
              opacity: isSelected ? 1 : 0
            }}
          >
            <IconDots size={15} />
          </ActionIcon>
        </Menu.Target>

        <Menu.Dropdown>
          <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={handleDelete}>
            Remove from Playlist
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
    </Group>
  )
}

export default PlaylistSongCard