import Playlist from '../../types/models/Playlist'
import { Center, Group, Stack, Text } from '@mantine/core'
import { useNavigate } from 'react-router-dom'
import { IconLayoutSidebarLeftExpand, IconPlaylist, IconTrash } from '@tabler/icons-react'
import { toast } from 'react-toastify'
import { useDeletePlaylistMutation } from '../../state/api/playlistsApi.ts'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { useDisclosure, useHover, useMergedRef } from '@mantine/hooks'
import { openPlaylistDrawer } from '../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../state/store.ts'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'
import PerfectRehearsalMenuItem from '../@ui/menu/item/PerfectRehearsalMenuItem.tsx'
import useDragSelectSelectable from '../../hooks/useDragSelectSelectable.ts'
import { MouseEvent } from 'react'
import SelectableAvatar from '../@ui/image/SelectableAvatar.tsx'

interface PlaylistCardProps {
  playlist: Playlist
}

function PlaylistCard({ playlist }: PlaylistCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const {
    ref: dragRef,
    isDragSelected,
    isDragSelecting
  } = useDragSelectSelectable<HTMLDivElement>(playlist.id)
  const { ref: hoverRef, hovered } = useHover<HTMLDivElement>()
  const ref = useMergedRef(dragRef, hoverRef)

  const [deletePlaylistMutation, { isLoading: isDeleteLoading }] = useDeletePlaylistMutation()

  const [openedMenu, { toggle: toggleMenu, close: closeMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const isSelected = openedMenu || hovered || isDragSelected

  function handleClick(e: MouseEvent) {
    if (e.ctrlKey || e.shiftKey) return
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
      aria-selected={isSelected}
      align={'center'}
      gap={0}
      style={{
        transition: '0.3s',
        ...(isSelected && { transform: 'scale(1.1)' })
      }}
    >
      <ContextMenu opened={openedMenu} onChange={toggleMenu} disabled={isDragSelecting}>
        <ContextMenu.Target>
          <SelectableAvatar
            ref={ref}
            id={playlist.id}
            radius={'10%'}
            w={'100%'}
            h={'unset'}
            src={playlist.imageUrl}
            alt={playlist.imageUrl && playlist.title}
            bg={'gray.5'}
            checkmarkSize={'28%'}
            isSelected={isDragSelected}
            sx={(theme) => ({
              aspectRatio: 1,
              cursor: 'pointer',
              transition: '0.3s',
              boxShadow: isSelected ? theme.shadows.xxl_hover : theme.shadows.xxl,
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
          </SelectableAvatar>
        </ContextMenu.Target>

        <ContextMenu.Dropdown>
          <ContextMenu.Item
            leftSection={<IconLayoutSidebarLeftExpand size={14} />}
            onClick={handleOpenDrawer}
          >
            Open Drawer
          </ContextMenu.Item>
          <PerfectRehearsalMenuItem id={playlist.id} closeMenu={closeMenu} type={'playlist'} />
          <ContextMenu.Divider />

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
