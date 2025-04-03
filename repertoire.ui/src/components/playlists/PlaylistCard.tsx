import Playlist from '../../types/models/Playlist'
import { Avatar, Center, Group, Menu, Stack, Text } from '@mantine/core'
import { useNavigate } from 'react-router-dom'
import { IconPlaylist, IconTrash } from '@tabler/icons-react'
import { toast } from 'react-toastify'
import { useDeletePlaylistMutation } from '../../state/api/playlistsApi.ts'
import useContextMenu from '../../hooks/useContextMenu.ts'
import { useState } from 'react'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { useDisclosure } from '@mantine/hooks'

interface PlaylistCardProps {
  playlist: Playlist
}

function PlaylistCard({ playlist }: PlaylistCardProps) {
  const navigate = useNavigate()

  const [deletePlaylistMutation, { isLoading: isDeleteLoading }] = useDeletePlaylistMutation()

  const [isImageHovered, setIsImageHovered] = useState(false)
  const [openedMenu, menuDropdownProps, { openMenu, closeMenu }] = useContextMenu()

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleClick() {
    navigate(`/playlist/${playlist.id}`)
  }

  async function handleDelete() {
    await deletePlaylistMutation(playlist.id).unwrap()
    toast.success(`${playlist.title} deleted!`)
  }

  return (
    <Stack
      aria-label={`playlist-card-${playlist.title}`}
      align={'center'}
      gap={0}
      style={{ transition: '0.3s', ...(isImageHovered && { transform: 'scale(1.1)' }) }}
    >
      <Menu shadow={'lg'} opened={openedMenu} onClose={closeMenu}>
        <Menu.Target>
          <Avatar
            onMouseEnter={() => setIsImageHovered(true)}
            onMouseLeave={() => setIsImageHovered(false)}
            radius={'10%'}
            w={'100%'}
            h={'unset'}
            src={playlist.imageUrl}
            alt={playlist.title}
            bg={'gray.5'}
            onClick={handleClick}
            onContextMenu={openMenu}
            sx={(theme) => ({
              aspectRatio: 1,
              cursor: 'pointer',
              transition: '0.3s',
              boxShadow: theme.shadows.xxl,
              '&:hover': { boxShadow: theme.shadows.xxl_hover }
            })}
          >
            <Center c={'white'}>
              <IconPlaylist size={'100%'} style={{ padding: '33%' }} />
            </Center>
          </Avatar>
        </Menu.Target>

        <Menu.Dropdown {...menuDropdownProps}>
          <Menu.Item c={'red'} leftSection={<IconTrash size={14} />} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>

      <Stack w={'100%'} pt={'xs'} style={{ overflow: 'hidden' }}>
        <Text fw={500} lineClamp={2} ta={'center'}>
          {playlist.title}
        </Text>
      </Stack>

      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={`Delete Playlist`}
        description={
          <Group gap={'xxs'}>
            <Text>Are you sure you want to delete</Text>
            <Text fw={600}>{playlist.title}</Text>
            <Text>?</Text>
          </Group>
        }
        onYes={handleDelete}
        isLoading={isDeleteLoading}
      />
    </Stack>
  )
}

export default PlaylistCard
