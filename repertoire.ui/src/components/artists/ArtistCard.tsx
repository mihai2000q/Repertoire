import Artist from '../../types/models/Artist.ts'
import { Center, Menu, Stack, Text } from '@mantine/core'
import { useNavigate } from 'react-router-dom'
import { IconLayoutSidebarLeftExpand, IconTrash } from '@tabler/icons-react'
import { useDisclosure, useHover, useMergedRef } from '@mantine/hooks'
import CustomIconUserAlt from '../@ui/icons/CustomIconUserAlt.tsx'
import { openArtistDrawer } from '../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../state/store.ts'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'
import DeleteArtistModal from '../@ui/modal/delete/DeleteArtistModal.tsx'
import PerfectRehearsalMenuItem from '../@ui/menu/item/PerfectRehearsalMenuItem.tsx'
import { MouseEvent } from 'react'
import SelectableAvatar from '../@ui/image/SelectableAvatar.tsx'
import useDragSelectSelectable from '../../hooks/useDragSelectSelectable.ts'

interface ArtistCardProps {
  artist: Artist
}

function ArtistCard({ artist }: ArtistCardProps) {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const {
    ref: dragRef,
    isDragSelected,
    isDragSelecting
  } = useDragSelectSelectable<HTMLDivElement>(artist.id)
  const { ref: hoverRef, hovered } = useHover<HTMLDivElement>()
  const ref = useMergedRef(dragRef, hoverRef)

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const isSelected = openedMenu || hovered || isDragSelected

  function handleClick(e: MouseEvent) {
    if (e.ctrlKey || e.shiftKey) return
    navigate(`/artist/${artist.id}`)
  }

  function handleOpenDrawer() {
    dispatch(openArtistDrawer(artist.id))
  }

  return (
    <Stack
      aria-label={`artist-card-${artist.name}`}
      align={'center'}
      gap={'xs'}
      style={{
        transition: '0.25s',
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
            id={artist.id}
            src={artist.imageUrl}
            alt={artist.imageUrl && artist.name}
            w={'100%'}
            h={'unset'}
            bg={'gray.0'}
            isSelected={isDragSelected}
            style={(theme) => ({
              aspectRatio: 1,
              cursor: 'pointer',
              transition: '0.3s',
              boxShadow: isSelected ? theme.shadows.xxl_hover : theme.shadows.xxl
            })}
            onClick={handleClick}
          >
            <Center c={'gray.7'}>
              <CustomIconUserAlt
                aria-label={`default-icon-${artist.name}`}
                size={'100%'}
                style={{ padding: '27%' }}
              />
            </Center>
          </SelectableAvatar>
        </ContextMenu.Target>

        <ContextMenu.Dropdown>
          <Menu.Item
            leftSection={<IconLayoutSidebarLeftExpand size={14} />}
            onClick={handleOpenDrawer}
          >
            Open Drawer
          </Menu.Item>
          <Menu.Divider />

          <AddToPlaylistMenuItem
            ids={[artist.id]}
            type={'artists'}
            closeMenu={closeMenu}
            disabled={artist.songsCount === 0}
          />
          <PerfectRehearsalMenuItem id={artist.id} closeMenu={closeMenu} type={'artist'} />
          <Menu.Divider />

          <Menu.Item c={'red'} leftSection={<IconTrash size={14} />} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>
      <Text px={'xs'} fw={600} ta={'center'} lineClamp={2}>
        {artist.name}
      </Text>

      <DeleteArtistModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        artist={artist}
        withName
      />
    </Stack>
  )
}

export default ArtistCard
