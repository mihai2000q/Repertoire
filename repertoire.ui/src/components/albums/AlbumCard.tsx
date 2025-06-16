import Album from '../../types/models/Album.ts'
import { Avatar, Center, Checkbox, Group, Stack, Text } from '@mantine/core'
import { useState } from 'react'
import { useAppDispatch } from '../../state/store.ts'
import { openAlbumDrawer, openArtistDrawer } from '../../state/slice/globalSlice.ts'
import { useNavigate } from 'react-router-dom'
import { IconLayoutSidebarLeftExpand, IconTrash, IconUser } from '@tabler/icons-react'
import { toast } from 'react-toastify'
import { useDeleteAlbumMutation } from '../../state/api/albumsApi.ts'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import CustomIconAlbumVinyl from '../@ui/icons/CustomIconAlbumVinyl.tsx'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'

interface AlbumCardProps {
  album: Album
}

function AlbumCard({ album }: AlbumCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [deleteAlbumMutation, { isLoading: isDeleteLoading }] = useDeleteAlbumMutation()

  const [deleteWithSongs, setDeleteWithSongs] = useState(false)

  const [isImageHovered, setIsImageHovered] = useState(false)
  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleClick() {
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

  async function handleDelete() {
    await deleteAlbumMutation({ id: album.id, withSongs: deleteWithSongs }).unwrap()
    toast.success(`${album.title} deleted!`)
  }

  return (
    <Stack
      aria-label={`album-card-${album.title}`}
      align={'center'}
      gap={0}
      style={{
        transition: '0.3s',
        ...((openedMenu || isImageHovered) && { transform: 'scale(1.1)' })
      }}
    >
      <ContextMenu shadow={'lg'} opened={openedMenu} onClose={closeMenu} onOpen={openMenu}>
        <ContextMenu.Target>
          <Avatar
            radius={'10%'}
            w={'100%'}
            h={'unset'}
            src={album.imageUrl}
            alt={album.imageUrl && album.title}
            bg={'gray.5'}
            sx={(theme) => ({
              aspectRatio: 1,
              cursor: 'pointer',
              transition: '0.3s',
              boxShadow: theme.shadows.xxl,
              '&:hover': { boxShadow: theme.shadows.xxl_hover },
              ...(openedMenu && { boxShadow: theme.shadows.xxl_hover })
            })}
            onMouseEnter={() => setIsImageHovered(true)}
            onMouseLeave={() => setIsImageHovered(false)}
            onClick={handleClick}
          >
            <Center c={'white'}>
              <CustomIconAlbumVinyl
                aria-label={`default-icon-${album.title}`}
                size={'100%'}
                style={{ padding: '37%' }}
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
            leftSection={<IconUser size={14} />}
            disabled={!album.artist}
            onClick={handleViewArtist}
          >
            View Artist
          </ContextMenu.Item>
          <AddToPlaylistMenuItem
            ids={[album.id]}
            type={'album'}
            closeMenu={closeMenu}
            disabled={album.songsCount === 0}
          />
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

      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={`Delete Album`}
        description={
          <Stack gap={5}>
            <Group gap={'xxs'}>
              <Text>Are you sure you want to delete</Text>
              <Text fw={600}>{album.title}</Text>
              <Text>?</Text>
            </Group>
            <Checkbox
              checked={deleteWithSongs}
              onChange={(event) => setDeleteWithSongs(event.currentTarget.checked)}
              label={'Delete all associated songs'}
              c={'dimmed'}
              styles={{ label: { paddingLeft: 8 } }}
            />
          </Stack>
        }
        onYes={handleDelete}
        isLoading={isDeleteLoading}
      />
    </Stack>
  )
}

export default AlbumCard
