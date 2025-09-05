import Album from '../../types/models/Album.ts'
import { Center, Stack, Text } from '@mantine/core'
import { useAppDispatch } from '../../state/store.ts'
import { openAlbumDrawer, openArtistDrawer } from '../../state/slice/globalSlice.ts'
import { useNavigate } from 'react-router-dom'
import { IconLayoutSidebarLeftExpand, IconTrash, IconUser } from '@tabler/icons-react'
import { useDisclosure, useHover, useMergedRef } from '@mantine/hooks'
import CustomIconAlbumVinyl from '../@ui/icons/CustomIconAlbumVinyl.tsx'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'
import DeleteAlbumModal from '../@ui/modal/delete/DeleteAlbumModal.tsx'
import PerfectRehearsalMenuItem from '../@ui/menu/item/PerfectRehearsalMenuItem.tsx'
import { MouseEvent } from 'react'
import useDragSelectSelectable from '../../hooks/useDragSelectSelectable.ts'
import SelectableAvatar from '../@ui/image/SelectableAvatar.tsx'

interface AlbumCardProps {
  album: Album
}

function AlbumCard({ album }: AlbumCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const {
    ref: dragRef,
    isDragSelected,
    isDragSelecting
  } = useDragSelectSelectable<HTMLDivElement>(album.id)
  const { ref: hoverRef, hovered } = useHover<HTMLDivElement>()
  const ref = useMergedRef(dragRef, hoverRef)

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const isSelected = openedMenu || hovered || isDragSelected

  function handleClick(e: MouseEvent) {
    if (e.ctrlKey || e.shiftKey) return
    navigate(`/album/${album.id}`)
  }

  function handleArtistClick() {
    dispatch(openArtistDrawer(album.artist.id))
  }

  function handleOpenDrawer() {
    dispatch(openAlbumDrawer(album.id))
  }

  function handleViewArtist() {
    navigate(`/artist/${album.artist.id}`)
  }

  return (
    <Stack
      aria-label={`album-card-${album.title}`}
      aria-selected={isSelected}
      align={'center'}
      gap={0}
      style={{
        transition: '0.3s',
        ...(isSelected && { transform: 'scale(1.1)' })
      }}
    >
      <ContextMenu
        opened={openedMenu}
        onClose={closeMenu}
        onOpen={openMenu}
        disabled={isDragSelecting}
      >
        <ContextMenu.Target>
          <SelectableAvatar
            ref={ref}
            id={album.id}
            radius={'10%'}
            w={'100%'}
            h={'unset'}
            src={album.imageUrl}
            alt={album.imageUrl && album.title}
            bg={'gray.5'}
            checkmarkSize={'28%'}
            isSelected={isDragSelected}
            sx={(theme) => ({
              aspectRatio: 1,
              cursor: 'pointer',
              transition: '0.3s',
              boxShadow: isSelected ? theme.shadows.xxl_hover : theme.shadows.xxl
            })}
            onClick={handleClick}
          >
            <Center c={'white'}>
              <CustomIconAlbumVinyl
                aria-label={`default-icon-${album.title}`}
                size={'100%'}
                style={{ padding: '37%' }}
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
          <ContextMenu.Item
            leftSection={<IconUser size={14} />}
            disabled={!album.artist}
            onClick={handleViewArtist}
          >
            View Artist
          </ContextMenu.Item>
          <ContextMenu.Divider />

          <AddToPlaylistMenuItem
            ids={[album.id]}
            type={'albums'}
            closeMenu={closeMenu}
            disabled={album.songsCount === 0}
          />
          <PerfectRehearsalMenuItem id={album.id} closeMenu={closeMenu} type={'album'} />
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

      <Stack w={'100%'} pt={'xs'} gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={600} lineClamp={2} ta={'center'}>
          {album.title}
        </Text>
        {album.artist ? (
          <Text
            fw={500}
            ta={'center'}
            c={'dimmed'}
            truncate={'end'}
            onClick={handleArtistClick}
            sx={{
              cursor: 'pointer',
              '&:hover': {
                textDecoration: 'underline'
              }
            }}
          >
            {album.artist.name}
          </Text>
        ) : (
          <Text c={'dimmed'} ta={'center'} fs={'oblique'}>
            Unknown
          </Text>
        )}
      </Stack>

      <DeleteAlbumModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        album={album}
        withName
      />
    </Stack>
  )
}

export default AlbumCard
