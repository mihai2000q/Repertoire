import Artist from '../../types/models/Artist.ts'
import { Avatar, Center, Checkbox, Group, Menu, Stack, Text } from '@mantine/core'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { IconLayoutSidebarLeftExpand, IconTrash } from '@tabler/icons-react'
import { toast } from 'react-toastify'
import useContextMenu from '../../hooks/useContextMenu.ts'
import { useDeleteArtistMutation } from '../../state/api/artistsApi.ts'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import CustomIconUserAlt from '../@ui/icons/CustomIconUserAlt.tsx'
import { openArtistDrawer } from '../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../state/store.ts'

interface ArtistCardProps {
  artist: Artist
}

function ArtistCard({ artist }: ArtistCardProps) {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()

  const [deleteArtistMutation, { isLoading: isDeleteLoading }] = useDeleteArtistMutation()

  const [deleteWithAssociations, setDeleteWithAssociations] = useState(false)

  const [isAvatarHovered, setIsAvatarHovered] = useState(false)

  const [openedMenu, menuDropdownProps, { openMenu, closeMenu }] = useContextMenu()

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleClick() {
    navigate(`/artist/${artist.id}`)
  }

  function handleOpenDrawer() {
    dispatch(openArtistDrawer(artist.id))
  }

  async function handleDelete() {
    await deleteArtistMutation({
      id: artist.id,
      withAlbums: deleteWithAssociations,
      withSongs: deleteWithAssociations
    }).unwrap()
    toast.success(`${artist.name} deleted!`)
  }

  return (
    <Stack
      aria-label={`artist-card-${artist.name}`}
      align={'center'}
      gap={'xs'}
      style={{
        transition: '0.25s',
        ...((openedMenu || isAvatarHovered) && { transform: 'scale(1.1)' })
      }}
    >
      <Menu shadow={'lg'} opened={openedMenu} onClose={closeMenu}>
        <Menu.Target>
          <Avatar
            onMouseEnter={() => setIsAvatarHovered(true)}
            onMouseLeave={() => setIsAvatarHovered(false)}
            src={artist.imageUrl}
            alt={artist.imageUrl && artist.name}
            w={'100%'}
            h={'unset'}
            bg={'gray.0'}
            style={(theme) => ({
              aspectRatio: 1,
              cursor: 'pointer',
              transition: '0.3s',
              boxShadow: openedMenu || isAvatarHovered ? theme.shadows.xxl_hover : theme.shadows.xxl
            })}
            onClick={handleClick}
            onContextMenu={openMenu}
          >
            <Center c={'gray.7'}>
              <CustomIconUserAlt
                aria-label={`default-icon-${artist.name}`}
                size={'100%'}
                style={{ padding: '27%' }}
              />
            </Center>
          </Avatar>
        </Menu.Target>

        <Menu.Dropdown {...menuDropdownProps}>
          <Menu.Item
            leftSection={<IconLayoutSidebarLeftExpand size={14} />}
            onClick={handleOpenDrawer}
          >
            Open Drawer
          </Menu.Item>
          <Menu.Item c={'red'} leftSection={<IconTrash size={14} />} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
      <Text px={'xs'} fw={600} ta={'center'} lineClamp={2}>
        {artist.name}
      </Text>

      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={'Delete Artist'}
        description={
          <Stack gap={'xs'}>
            <Group gap={'xxs'}>
              <Text>Are you sure you want to delete</Text>
              <Text fw={600}>{artist.name}</Text>
              <Text>?</Text>
            </Group>
            <Checkbox
              checked={deleteWithAssociations}
              onChange={(event) => setDeleteWithAssociations(event.currentTarget.checked)}
              label={
                <Text c={'dimmed'}>
                  Delete all associated <b>albums</b> and <b>songs</b>
                </Text>
              }
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

export default ArtistCard
