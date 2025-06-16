import Playlist from '../../types/models/Playlist'
import { Avatar, Center, Group, Stack, Text } from '@mantine/core'
import { useNavigate } from 'react-router-dom'
import { IconLayoutSidebarLeftExpand, IconPlaylist, IconTrash } from '@tabler/icons-react'
import { toast } from 'react-toastify'
import { useDeletePlaylistMutation } from '../../state/api/playlistsApi.ts'
import { useState } from 'react'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { openPlaylistDrawer } from '../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../state/store.ts'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'

interface PlaylistCardProps {
  playlist: Playlist
}

function PlaylistCard({ playlist }: PlaylistCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [deletePlaylistMutation, { isLoading: isDeleteLoading }] = useDeletePlaylistMutation()

  const [isImageHovered, setIsImageHovered] = useState(false)
  const [openedMenu, { toggle: toggleMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleClick() {
    navigate(`/playlist/${playlist.id}`)
  }

  function handleOpenDrawer() {
    dispatch(openPlaylistDrawer(playlist.id))
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
      style={{
        transition: '0.3s',
        ...((openedMenu || isImageHovered) && { transform: 'scale(1.1)' })
      }}
    >
      <ContextMenu shadow={'lg'} opened={openedMenu} onChange={toggleMenu}>
        <ContextMenu.Target>
          <Avatar
            onMouseEnter={() => setIsImageHovered(true)}
            onMouseLeave={() => setIsImageHovered(false)}
            radius={'10%'}
            w={'100%'}
            h={'unset'}
            src={playlist.imageUrl}
            alt={playlist.imageUrl && playlist.title}
            bg={'gray.5'}
            sx={(theme) => ({
              aspectRatio: 1,
              cursor: 'pointer',
              transition: '0.3s',
              boxShadow: theme.shadows.xxl,
              '&:hover': { boxShadow: theme.shadows.xxl_hover },
              ...(openedMenu && { boxShadow: theme.shadows.xxl_hover })
            })}
            onClick={handleClick}
          >
            <Center c={'white'}>
              <IconPlaylist
                aria-label={`default-icon-${playlist.title}`}
                size={'100%'}
                style={{ padding: '33%' }}
              />
            </Center>
          </Avatar>
        </ContextMenu.Target>

        <ContextMenu.Dropdown>
          <ContextMenu.Item
            leftSection={<IconLayoutSidebarLeftExpand size={14} />}
            onClick={handleOpenDrawer}
          >
            Open Drawer
          </ContextMenu.Item>
          <ContextMenu.Item
            c={'red'}
            leftSection={<IconTrash size={14} />}
            onClick={openDeleteWarning}
          >
            Delete
          </ContextMenu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>

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
